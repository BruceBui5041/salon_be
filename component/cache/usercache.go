package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"video_server/appconst"
	"video_server/common"
	"video_server/component/logger"
	models "video_server/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type CacheUser struct {
	Status            string              `json:"status"`
	LastName          string              `json:"lastname" gorm:"column:lastname;"`
	FirstName         string              `json:"firstname" gorm:"column:firstname;"`
	Email             string              `json:"email"`
	ProfilePictureURL string              `json:"profile_picture_url"`
	Roles             []CacheRoleRes      `json:"roles"`
	Enrollments       []CacheEnrollment   `json:"enrollments"`
	Auths             []models.UserAuth   `json:"auths"`
	FakeId            *common.UID         `json:"id"`
	UserProfile       *models.UserProfile `json:"user_profile,omitempty"`
}

type CacheRoleRes struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description"`
}

type CacheEnrollment struct {
	common.SQLModel `json:",inline"`
	EnrolledAt      time.Time   `json:"enrolled_at"`
	Course          CacheCourse `json:"course"`
}

type CacheCourse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
}

func (ac *appCache) SetUserCache(ctx context.Context, user *models.User) error {
	user.Mask(false)
	cacheKey := fmt.Sprintf("%s:%s", appconst.UserPrefix, user.GetFakeId())
	var cacheUser CacheUser
	err := copier.Copy(&cacheUser, user)
	if err != nil {
		return err
	}
	cacheUser.FakeId = user.FakeId

	// Set in DynamoDB
	return ac.setUserCacheDynamoDB(ctx, cacheKey, cacheUser)
}

func (ac *appCache) GetUserCache(ctx context.Context, fakeUserId string) (string, error) {
	start := time.Now()
	key := fmt.Sprintf("%s:%s", appconst.UserPrefix, fakeUserId)

	// Try to get from memory cache first
	if cachedUser, found := ac.GetLocalCache().Get(key); found {
		cacheUserJson, err := json.Marshal(cachedUser)
		if err == nil {
			return string(cacheUserJson), nil
		}
		// If marshaling fails, log the error and continue to DynamoDB
		logger.AppLogger.Warn(ctx, "Failed to marshal cached user", zap.Error(err))
	}

	// If not found in memory cache, try DynamoDB
	userJson, err := ac.getUserCacheDynamoDB(ctx, key, fakeUserId)
	if err != nil {
		return "", err
	}

	// If found in DynamoDB, set it in the local cache
	if userJson != "" {
		var cacheUser CacheUser
		err = json.Unmarshal([]byte(userJson), &cacheUser)
		if err == nil {
			userCacheTTL := viper.GetInt("DYNAMODB_USER_CACHE_TTL")
			ac.GetLocalCache().Set(key, cacheUser, time.Duration(userCacheTTL)*time.Hour)
		} else {
			logger.AppLogger.Warn(ctx, "Failed to unmarshal user from DynamoDB", zap.Error(err))
		}
	}

	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info(ctx, "Cache Get duration",
			zap.String("key", key),
			zap.Float64("duration_ms", durationMs))
	}()

	return userJson, nil
}

func (ac *appCache) setUserCacheDynamoDB(ctx context.Context, cacheKey string, cacheUser CacheUser) error {
	cacheUserJson, err := json.Marshal(cacheUser)
	if err != nil {
		logger.AppLogger.Warn(ctx, "cannot parse user for cache", zap.Error(err))
		return err
	}

	userCacheTTL := viper.GetInt("DYNAMODB_USER_CACHE_TTL")
	ttl := time.Now().Add(time.Duration(userCacheTTL) * time.Hour).Unix()
	item := map[string]*dynamodb.AttributeValue{
		"cachekey": {
			S: aws.String(cacheKey),
		},
		"userid": {
			S: aws.String(cacheUser.FakeId.String()),
		},
		"value": {
			S: aws.String(string(cacheUserJson)),
		},
		"ttl": {
			N: aws.String(strconv.FormatInt(ttl, 10)),
		},
	}

	cacheDuration := time.Duration(userCacheTTL) * time.Hour

	return ac.putItemDynamoDB(
		ctx,
		os.Getenv("DYNAMODB_USER_TABLE_NAME"),
		item,
		cacheKey,
		cacheUser,
		cacheDuration,
	)
}

func (ac *appCache) getUserCacheDynamoDB(ctx context.Context, key string, fakeUserId string) (string, error) {
	userTableName := os.Getenv("DYNAMODB_USER_TABLE_NAME")
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(key),
			},
			"userid": {
				S: aws.String(fakeUserId),
			},
		},
		TableName:      aws.String(userTableName),
		ConsistentRead: aws.Bool(true),
	}
	result, err := ac.GetDynamoDB().GetItem(input)
	if err != nil {
		return "", err
	}
	if result.Item == nil {
		return "", nil
	}
	value, ok := result.Item["value"]
	if !ok || value.S == nil {
		logger.AppLogger.Warn(ctx, "dynamoDB not found", zap.String("key", key))
		return "", nil
	}
	return *value.S, nil
}

func (ac *appCache) DeleteUserCache(ctx context.Context, userId string) error {
	userTableName := os.Getenv("DYNAMODB_USER_TABLE_NAME")
	cacheKey := fmt.Sprintf("%s:%s", appconst.UserPrefix, userId)

	// Delete from memory cache
	ac.GetLocalCache().Delete(cacheKey)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(cacheKey),
			},
			"userid": {
				S: aws.String(userId),
			},
		},
		TableName: aws.String(userTableName),
	}

	if err := ac.deleteItemFromDynamoDB(ctx, input); err != nil {
		return err
	}

	logger.AppLogger.Info(ctx, "User cache deleted successfully", zap.String("userId", userId))
	return nil
}
