package cache

import (
	"context"
	"fmt"
	"os"
	"salon_be/appconst"
	"salon_be/component/logger"
	models "salon_be/model"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type VideoCacheInfo struct {
	VideoID    string
	VideoURL   string
	CourseSlug string
	ExpiresAt  int64
}

func (ac *appCache) SetVideoCache(ctx context.Context, courseSlug string, video models.Video) error {
	video.Mask(false)

	cacheKey := fmt.Sprintf("%s:%s:%s", appconst.VideoURLPrefix, courseSlug, video.GetFakeId())

	videoTableName := os.Getenv("DYNAMODB_VIDEO_TABLE_NAME")
	videoCacheTTL := viper.GetInt("DYNAMODB_VIDEO_CACHE_TTL")

	var expirationTime int64
	var cacheDuration time.Duration
	if videoCacheTTL > 0 {
		cacheDuration = time.Duration(videoCacheTTL) * time.Hour
		expirationTime = time.Now().Add(cacheDuration).Unix()
	} else {
		cacheDuration = cache.NoExpiration
		expirationTime = 0 // 0 can represent no expiration in DynamoDB
	}

	videoCacheInfo := VideoCacheInfo{
		VideoID:    video.GetFakeId(),
		VideoURL:   video.VideoURL,
		CourseSlug: courseSlug,
		ExpiresAt:  expirationTime,
	}

	// Prepare the item for DynamoDB
	item := map[string]*dynamodb.AttributeValue{
		"cachekey": {
			S: aws.String(cacheKey),
		},
		"videoid": {
			S: aws.String(videoCacheInfo.VideoID),
		},
		"videourl": {
			S: aws.String(videoCacheInfo.VideoURL),
		},
		"courseslug": {
			S: aws.String(videoCacheInfo.CourseSlug),
		},
		"ttl": {
			N: aws.String(strconv.FormatInt(videoCacheInfo.ExpiresAt, 10)),
		},
	}

	// Use the putItemDynamoDB function
	err := ac.putItemDynamoDB(
		ctx,
		videoTableName,
		item,
		cacheKey,
		videoCacheInfo,
		cacheDuration,
	)

	if err != nil {
		logger.AppLogger.Warn(ctx, "Failed to set video cache",
			zap.Error(err),
			zap.String("courseSlug", courseSlug),
			zap.Uint32("videoId", video.Id),
			zap.Int("cacheTTL", videoCacheTTL))
	}

	return err
}

func (ac *appCache) GetVideoCache(ctx context.Context, courseSlug string, videoId string) (*VideoCacheInfo, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", appconst.VideoURLPrefix, courseSlug, videoId)

	// Try to get from memory cache first
	if cachedVideoInfo, found := ac.GetLocalCache().Get(cacheKey); found {
		if videoInfo, ok := cachedVideoInfo.(*VideoCacheInfo); ok {
			return videoInfo, nil
		}
		// If type assertion fails, log the error and continue to DynamoDB
		logger.AppLogger.Warn(ctx, "Failed to retrieve video info from memory cache", zap.String("key", cacheKey))
	}

	// If not found in memory cache, try DynamoDB
	videoInfo, err := ac.getVideoCacheDynamoDB(ctx, cacheKey, videoId)
	if err != nil {
		return nil, err
	}

	// If found in DynamoDB, set it in the local cache
	if videoInfo != nil {
		videoCacheTTL := viper.GetInt("DYNAMODB_VIDEO_CACHE_TTL")
		ac.GetLocalCache().Set(cacheKey, videoInfo, time.Duration(videoCacheTTL)*time.Hour)
	}

	return videoInfo, nil
}

func (ac *appCache) getVideoCacheDynamoDB(ctx context.Context, cacheKey string, videoId string) (*VideoCacheInfo, error) {
	start := time.Now()
	videoTableName := os.Getenv("DYNAMODB_VIDEO_TABLE_NAME")

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"cachekey": {
				S: aws.String(cacheKey),
			},
			"videoid": {
				S: aws.String(videoId),
			},
		},
		TableName:      aws.String(videoTableName),
		ConsistentRead: aws.Bool(false),
	}

	result, err := ac.GetDynamoDB().GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	videoInfo := &VideoCacheInfo{}

	if value, ok := result.Item["videourl"]; ok && value.S != nil {
		videoInfo.VideoURL = *value.S
	}

	if value, ok := result.Item["videoid"]; ok && value.S != nil {
		videoInfo.VideoID = *value.S
	}

	if value, ok := result.Item["courseslug"]; ok && value.S != nil {
		videoInfo.CourseSlug = *value.S
	}

	if value, ok := result.Item["ttl"]; ok && value.N != nil {
		if expiresAt, err := strconv.ParseInt(*value.N, 10, 64); err == nil {
			videoInfo.ExpiresAt = expiresAt
		}
	}

	if videoInfo.VideoURL == "" {
		logger.AppLogger.Warn(ctx, "dynamoDB not found", zap.String("key", cacheKey))
		return nil, nil
	}
	defer func() {
		duration := time.Since(start)
		durationMs := float64(duration) / float64(time.Millisecond)
		logger.AppLogger.Info(ctx, "Cache Get duration from dynamoDB",
			zap.String("key", cacheKey),
			zap.Float64("duration_ms", durationMs))
	}()
	return videoInfo, nil
}
