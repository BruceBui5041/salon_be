package component

import (
	"context"
	"salon_be/component/cache"
	models "salon_be/model"
	pb "salon_be/proto/video_service/video_service"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetLocalPubSub() LocalPubSub
	GetVideoProcessingClient() pb.VideoProcessingServiceClient
	SecretKey() string
	GetAppCache() AppCache
	GetAppQueue() AppQueue
	GetS3Client() *s3.S3
	GetCronJob() CronJob
}

type AppQueue interface {
	StartSQSMessageListener(
		ctx context.Context,
		appContext AppContext,
		consumeTopics []string,
		processMessageHandler func(ctx context.Context, appContext AppContext, msg *sqs.Message) error,
	)
	SendSQSMessage(ctx context.Context, topic, groupId, messageBody string) error
	ReceiveSQSMessage(ctx context.Context, queueURL string) ([]*sqs.Message, error)
	PollSQSMessages(ctx context.Context, appContext AppContext, queueURL string, handler func(context.Context, AppContext, *sqs.Message) error)
	CreateSQSQueues(ctx context.Context) map[string]*string
}

type CronJob interface {
	Start()
	RegisterVideoJobs(ctx context.Context, appCtx AppContext) error
	RegisterCourseJobs(ctx context.Context, appCtx AppContext) error
}

type AppCache interface {
	GetUserCache(ctx context.Context, fakeUserId string) (string, error)
	SetUserCache(ctx context.Context, user *models.User) error
	DeleteUserCache(ctx context.Context, userId string) error

	SetVideoCache(ctx context.Context, courseSlug string, video models.Video) error
	GetVideoCache(ctx context.Context, courseSlug string, videoId string) (*cache.VideoCacheInfo, error)

	SetEnrollmentCache(ctx context.Context, enrollment *cache.EnrollmentCache) error
	GetEnrollmentCache(ctx context.Context, courseSlug, userId string) (*cache.EnrollmentCache, error)
	DeleteEnrollmentCache(ctx context.Context, courseSlug, userId string) error

	CreateDynamoDBTables(ctx context.Context, tables []cache.DynamoDBTableDefinition) error
	GetDynamoDBTableDefinitions() []cache.DynamoDBTableDefinition
}

type LocalPubSub interface {
	GetUnblockPubSub() *gochannel.GoChannel
	// The publisher of this PubSub will be block util its message got acked
	GetBlockPubSub() *gochannel.GoChannel
}

type appCtx struct {
	db                 *gorm.DB
	localPubSub        LocalPubSub
	videoServiceClient pb.VideoProcessingServiceClient
	jwtSecretKey       string
	appCache           AppCache
	awsSession         *session.Session
	s3Client           *s3.S3
	appqueue           AppQueue
	cron               CronJob
}

func NewAppContext(
	db *gorm.DB,
	localPubSub LocalPubSub,
	videoServiceClient pb.VideoProcessingServiceClient,
	jwtSecretKey string,
	appCache AppCache,
	awsSession *session.Session,
	appqueue AppQueue,
	cron CronJob,
	s3Client *s3.S3,
) *appCtx {

	return &appCtx{
		db:                 db,
		localPubSub:        localPubSub,
		videoServiceClient: videoServiceClient,
		jwtSecretKey:       jwtSecretKey,
		appCache:           appCache,
		awsSession:         awsSession,
		appqueue:           appqueue,
		cron:               cron,
		s3Client:           s3Client,
	}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}

func (ctx *appCtx) GetLocalPubSub() LocalPubSub {
	return ctx.localPubSub
}

func (ctx *appCtx) GetVideoProcessingClient() pb.VideoProcessingServiceClient {
	return ctx.videoServiceClient
}

func (ctx *appCtx) SecretKey() string {
	return ctx.jwtSecretKey
}

func (ctx *appCtx) GetAppCache() AppCache {
	return ctx.appCache
}

func (ctx *appCtx) GetS3Client() *s3.S3 {
	return ctx.s3Client
}

func (ctx *appCtx) GetAppQueue() AppQueue {
	return ctx.appqueue
}

func (ctx *appCtx) GetCronJob() CronJob {
	return ctx.cron
}
