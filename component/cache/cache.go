package cache

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/patrickmn/go-cache"
)

type appCache struct {
	localCache *cache.Cache
	dynamodb   *dynamodb.DynamoDB
}

func NewAppCache(
	sess *session.Session,
) (*appCache, error) {
	client := dynamodb.New(sess)
	localCache := cache.New(1*time.Hour, 2*time.Hour)
	return &appCache{
		dynamodb:   client,
		localCache: localCache,
	}, nil
}

func (ac *appCache) GetLocalCache() *cache.Cache {
	return ac.localCache
}
