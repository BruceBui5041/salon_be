package cache

import (
	"context"
	"fmt"
	"salon_be/component/logger"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
)

type DynamoDBTableDefinition struct {
	Name                 string
	AttributeDefinitions []*dynamodb.AttributeDefinition
	KeySchema            []*dynamodb.KeySchemaElement
	ReadCapacity         int64
	WriteCapacity        int64
	TTLAttributeName     string
}

func (ac *appCache) GetDynamoDB() *dynamodb.DynamoDB {
	return ac.dynamodb
}

func (ac *appCache) CreateDynamoDBTables(ctx context.Context, tables []DynamoDBTableDefinition) error {
	for _, table := range tables {
		// Check if the table already exists
		exists, err := ac.tableExists(ctx, table.Name)
		if err != nil {
			logger.AppLogger.Error(ctx, "Error checking if table exists",
				zap.String("table", table.Name),
				zap.Error(err))
			return err
		}

		if exists {
			logger.AppLogger.Info(ctx, "Table already exists",
				zap.String("table", table.Name))
		} else {
			// Table doesn't exist, so create it
			input := &dynamodb.CreateTableInput{
				AttributeDefinitions: table.AttributeDefinitions,
				KeySchema:            table.KeySchema,
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(table.ReadCapacity),
					WriteCapacityUnits: aws.Int64(table.WriteCapacity),
				},
				TableName: aws.String(table.Name),
			}

			_, err = ac.dynamodb.CreateTable(input)
			if err != nil {
				logger.AppLogger.Error(ctx, "Error creating table",
					zap.String("table", table.Name),
					zap.Error(err))
				return err
			}

			logger.AppLogger.Info(ctx, "Table created successfully",
				zap.String("table", table.Name))
		}

		// Wait for the table to be active
		err = ac.waitForTableActive(ctx, table.Name)
		if err != nil {
			logger.AppLogger.Error(ctx, "Error waiting for table to become active",
				zap.String("table", table.Name),
				zap.Error(err))
			return err
		}

		// Enable TTL for the table
		if table.TTLAttributeName != "" {
			err = ac.enableTTLForTable(ctx, table.Name, table.TTLAttributeName)
			if err != nil {
				logger.AppLogger.Error(ctx, "Error enabling TTL for table",
					zap.String("table", table.Name),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}
func (ac *appCache) waitForTableActive(ctx context.Context, tableName string) error {
	logger.AppLogger.Info(ctx, "Waiting for table to become active",
		zap.String("table", tableName))

	for i := 0; i < 60; i++ { // Wait for up to 5 minutes (60 * 5 seconds)
		input := &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}
		result, err := ac.dynamodb.DescribeTable(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeResourceNotFoundException:
					// Table doesn't exist yet, wait and retry
					time.Sleep(5 * time.Second)
					continue
				default:
					return err
				}
			}
			return err
		}

		if *result.Table.TableStatus == dynamodb.TableStatusActive {
			logger.AppLogger.Info(ctx, "Table is now active",
				zap.String("table", tableName))
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for table %s to become active", tableName)
}

func (ac *appCache) enableTTLForTable(ctx context.Context, tableName, ttlAttributeName string) error {
	// First, check if TTL is already enabled
	describeInput := &dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(tableName),
	}

	describeOutput, err := ac.dynamodb.DescribeTimeToLive(describeInput)
	if err != nil {
		logger.AppLogger.Error(ctx, "Error describing TTL for table",
			zap.String("table", tableName),
			zap.Error(err))
		return err
	}

	if describeOutput.TimeToLiveDescription != nil &&
		describeOutput.TimeToLiveDescription.TimeToLiveStatus != nil &&
		*describeOutput.TimeToLiveDescription.TimeToLiveStatus == dynamodb.TimeToLiveStatusEnabled {
		logger.AppLogger.Info(ctx, "TTL is already enabled for table",
			zap.String("table", tableName),
			zap.String("ttlAttribute", *describeOutput.TimeToLiveDescription.AttributeName))
		return nil
	}

	// If TTL is not enabled, proceed with enabling it
	input := &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String(ttlAttributeName),
			Enabled:       aws.Bool(true),
		},
	}

	_, err = ac.dynamodb.UpdateTimeToLive(input)
	if err != nil {
		// Check if the error is because TTL is already enabled
		if strings.Contains(err.Error(), "TimeToLive is already enabled") {
			logger.AppLogger.Info(ctx, "TTL is already enabled for table",
				zap.String("table", tableName),
				zap.String("ttlAttribute", ttlAttributeName))
			return nil
		}
		logger.AppLogger.Error(ctx, "Error updating TTL for table",
			zap.String("table", tableName),
			zap.String("ttlAttribute", ttlAttributeName),
			zap.Error(err))
		return err
	}

	logger.AppLogger.Info(ctx, "TTL enabled for table",
		zap.String("table", tableName),
		zap.String("ttlAttribute", ttlAttributeName))
	return nil
}

func (ac *appCache) tableExists(ctx context.Context, tableName string) (bool, error) {
	input := &dynamodb.ListTablesInput{}
	for {
		result, err := ac.dynamodb.ListTables(input)
		if err != nil {
			logger.AppLogger.Error(ctx, "Error listing tables",
				zap.Error(err))
			return false, err
		}

		for _, n := range result.TableNames {
			if *n == tableName {
				return true, nil
			}
		}

		input.ExclusiveStartTableName = result.LastEvaluatedTableName
		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return false, nil
}

// putItemDynamoDB puts an item into DynamoDB and sets it in the local cache
func (ac *appCache) putItemDynamoDB(
	ctx context.Context,
	tableName string,
	item map[string]*dynamodb.AttributeValue,
	cacheKey string,
	cacheValue interface{},
	cacheDuration time.Duration,
) error {
	// Put item into DynamoDB
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err := ac.dynamodb.PutItem(input)
	if err != nil {
		logger.AppLogger.Error(ctx, "Error putting item in DynamoDB",
			zap.String("table", tableName),
			zap.Error(err))
		return err
	}

	// Set item in local cache
	ac.localCache.Set(cacheKey, cacheValue, cacheDuration)

	logger.AppLogger.Debug(ctx, "Item put in DynamoDB and local cache",
		zap.String("table", tableName),
		zap.String("cacheKey", cacheKey),
		zap.Duration("cacheDuration", cacheDuration))

	return nil
}

func (ac *appCache) deleteItemFromDynamoDB(ctx context.Context, input *dynamodb.DeleteItemInput) error {
	_, err := ac.GetDynamoDB().DeleteItem(input)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to delete item from DynamoDB",
			zap.String("tableName", *input.TableName),
			zap.Any("key", input.Key),
			zap.Error(err))
		return err
	}

	logger.AppLogger.Info(ctx, "Item deleted successfully from DynamoDB",
		zap.String("tableName", *input.TableName),
		zap.Any("key", input.Key))
	return nil
}
