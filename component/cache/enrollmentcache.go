package cache

import (
	"context"
	"fmt"
	"os"
	"salon_be/appconst"
	"salon_be/component/logger"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type EnrollmentCache struct {
	UserId            string `json:"user_id"`
	ServiceId         string `json:"service_id"`
	ServiceSlug       string `json:"service_slug"`
	EnrollmentId      string `json:"enrollment_id"`
	PaymentId         string `json:"payment_id"`
	TransactionStatus string `json:"transaction_status"`
}

func (ac *appCache) SetEnrollmentCache(ctx context.Context, enrollmentCache *EnrollmentCache) error {
	cacheKey := fmt.Sprintf(
		"%s:%s:%s",
		appconst.EnrollmentPrefix,
		enrollmentCache.ServiceSlug,
		enrollmentCache.UserId,
	)

	enrollmentTableName := os.Getenv("DYNAMODB_ENROLLMENT_TABLE_NAME")
	enrollmentCacheTTL := viper.GetInt("DYNAMODB_ENROLLMENT_CACHE_TTL")

	var expirationTime int64
	var cacheDuration time.Duration
	if enrollmentCacheTTL > 0 {
		cacheDuration = time.Duration(enrollmentCacheTTL) * time.Hour
		expirationTime = time.Now().Add(cacheDuration).Unix()
	} else {
		cacheDuration = cache.NoExpiration
		expirationTime = 0
	}

	item := map[string]*dynamodb.AttributeValue{
		"cachekey": {
			S: aws.String(cacheKey),
		},
		"userid": {
			S: aws.String(enrollmentCache.UserId),
		},
		"serviceslug": {
			S: aws.String(enrollmentCache.ServiceSlug),
		},
		"serviceid": {
			S: aws.String(enrollmentCache.ServiceId),
		},
		"enrollmentid": {
			S: aws.String(enrollmentCache.EnrollmentId),
		},
		"paymentid": {
			S: aws.String(enrollmentCache.PaymentId),
		},
		"transactionstatus": {
			S: aws.String(enrollmentCache.TransactionStatus),
		},
		"ttl": {
			N: aws.String(fmt.Sprintf("%d", expirationTime)),
		},
	}

	err := ac.putItemDynamoDB(
		ctx,
		enrollmentTableName,
		item,
		cacheKey,
		fmt.Sprintf("%s:%s", enrollmentCache.ServiceSlug, enrollmentCache.UserId),
		cacheDuration,
	)

	if err != nil {
		logger.AppLogger.Warn(ctx, "Failed to set enrollment cache",
			zap.Error(err),
			zap.String("enrollmentId", enrollmentCache.EnrollmentId),
			zap.String("userId", enrollmentCache.UserId),
			zap.String("serviceId", enrollmentCache.ServiceId),
			zap.String("serviceSlug", enrollmentCache.ServiceSlug),
			zap.Int("cacheTTL", enrollmentCacheTTL))
	}

	return err
}

func (ac *appCache) GetEnrollmentCache(ctx context.Context, serviceSlug, userId string) (*EnrollmentCache, error) {
	start := time.Now()
	cacheKey := fmt.Sprintf("%s:%s:%s", appconst.EnrollmentPrefix, serviceSlug, userId)

	if cachedEnrollment, found := ac.GetLocalCache().Get(cacheKey); found {
		if enrollment, ok := cachedEnrollment.(*EnrollmentCache); ok {
			return enrollment, nil
		}
		logger.AppLogger.Warn(ctx, "Failed to retrieve enrollment from memory cache", zap.String("key", cacheKey))
	}

	enrollmentCache, err := ac.getEnrollmentCacheDynamoDB(ctx, cacheKey, serviceSlug)
	if err != nil {
		return nil, err
	}

	enrollmentCacheTTL := viper.GetInt("DYNAMODB_ENROLLMENT_CACHE_TTL")
	ac.GetLocalCache().Set(cacheKey, enrollmentCache, time.Duration(enrollmentCacheTTL)*time.Hour)

	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info(ctx, "Cache Get duration from dynamoDB",
			zap.String("key", cacheKey),
			zap.Float64("duration_ms", durationMs))
	}()

	return enrollmentCache, nil
}

func (ac *appCache) getEnrollmentCacheDynamoDB(ctx context.Context, cacheKey, serviceSlug string) (*EnrollmentCache, error) {
	enrollmentTableName := os.Getenv("DYNAMODB_ENROLLMENT_TABLE_NAME")

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(cacheKey),
			},
			"serviceslug": {
				S: aws.String(serviceSlug),
			},
		},
		TableName:      aws.String(enrollmentTableName),
		ConsistentRead: aws.Bool(false),
	}

	result, err := ac.GetDynamoDB().GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	enrollmentCache := &EnrollmentCache{}

	if v, ok := result.Item["userid"]; ok && v.S != nil {
		enrollmentCache.UserId = *v.S
	}
	if v, ok := result.Item["serviceid"]; ok && v.S != nil {
		enrollmentCache.ServiceId = *v.S
	}
	if v, ok := result.Item["serviceslug"]; ok && v.S != nil {
		enrollmentCache.ServiceSlug = *v.S
	}
	if v, ok := result.Item["enrollmentid"]; ok && v.S != nil {
		enrollmentCache.EnrollmentId = *v.S
	}
	if v, ok := result.Item["paymentid"]; ok && v.S != nil {
		enrollmentCache.PaymentId = *v.S
	}
	if v, ok := result.Item["transactionstatus"]; ok && v.S != nil {
		enrollmentCache.TransactionStatus = *v.S
	}

	if enrollmentCache.UserId == "" || enrollmentCache.ServiceSlug == "" {
		logger.AppLogger.Warn(ctx, "dynamoDB item incomplete", zap.String("key", cacheKey))
		return nil, nil
	}

	return enrollmentCache, nil
}

func (ac *appCache) DeleteEnrollmentCache(ctx context.Context, serviceSlug, userId string) error {
	cacheKey := fmt.Sprintf("%s:%s:%s", appconst.EnrollmentPrefix, serviceSlug, userId)

	ac.GetLocalCache().Delete(cacheKey)

	logger.AppLogger.Info(ctx, "Enrollment cache deleted successfully", zap.String("serviceSlug", serviceSlug), zap.String("userId", userId))
	return nil
}
