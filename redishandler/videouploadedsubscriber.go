package redishandler

import (
	"context"
	"encoding/json"
	"log"
	"video_server/appconst"
	"video_server/watermill/messagemodel"

	"github.com/go-redis/redis/v8"
)

func StartRedisSubscribers(redisClient *redis.Client) {
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, appconst.TopicNewVideoUploaded)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Received message from Redis channel %s: %s", msg.Channel, msg.Payload)

		// Parse the message payload
		var videoInfo *messagemodel.RequestProcessVideoInfo
		err := json.Unmarshal([]byte(msg.Payload), &videoInfo)
		if err != nil {
			log.Printf("Error parsing message payload: %v", err)
			continue
		}

		// go watermill.PublishVideoUploadedEvent(videoInfo)
	}
}
