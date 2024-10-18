package watermill

import (
	"context"
	"log"
	"video_server/appconst"
	"video_server/component"

	"github.com/ThreeDotsLabs/watermill/message"
)

func StartSubscribers(appCtx component.AppContext) {
	// Define topics and their handlers
	topicHandlers := map[string]MessageHandler{
		appconst.TopicNewVideoUploaded: HandleNewVideoUpload,
		appconst.TopicVideoProcessed:   HandleVideoProcessed,
		appconst.TopicEnrollmentChange: EnrollmentChangeHandler,
		appconst.TopicUserUpdated:      UserUpdatedHandler,
		appconst.TopicReceivedWSMsg:    ReceivedWSMsgHandler,
		// Add more topics and handlers as needed
	}

	// Subscribe to all topics
	messageUnBlockPubChannels := make(map[string]<-chan *message.Message)
	for topic := range topicHandlers {
		messages, err := appCtx.GetLocalPubSub().GetUnblockPubSub().Subscribe(context.Background(), topic)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
		}
		messageUnBlockPubChannels[topic] = messages
	}

	// Process messages from all topics
	go processMessages(appCtx, messageUnBlockPubChannels, topicHandlers)

	messageBlockPubChannels := make(map[string]<-chan *message.Message)
	for topic := range topicHandlers {
		messages, err := appCtx.GetLocalPubSub().GetBlockPubSub().Subscribe(context.Background(), topic)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
		}
		messageBlockPubChannels[topic] = messages
	}

	// Process messages from all topics
	go processMessages(appCtx, messageBlockPubChannels, topicHandlers)

}

func processMessages(
	appCtx component.AppContext,
	messageChannels map[string]<-chan *message.Message,
	topicHandlers map[string]MessageHandler,
) {
	for {
		select {
		case msg := <-messageChannels[appconst.TopicNewVideoUploaded]:
			topicHandlers[appconst.TopicNewVideoUploaded](appCtx, msg)
		case msg := <-messageChannels[appconst.TopicVideoProcessed]:
			topicHandlers[appconst.TopicVideoProcessed](appCtx, msg)
		case msg := <-messageChannels[appconst.TopicEnrollmentChange]:
			topicHandlers[appconst.TopicEnrollmentChange](appCtx, msg)
		case msg := <-messageChannels[appconst.TopicUserUpdated]:
			topicHandlers[appconst.TopicUserUpdated](appCtx, msg)
		case msg := <-messageChannels[appconst.TopicReceivedWSMsg]:
			topicHandlers[appconst.TopicReceivedWSMsg](appCtx, msg)
			// Add more cases for additional topics
		}
	}
}
