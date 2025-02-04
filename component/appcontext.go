package component

import (
	"context"
	"salon_be/component/cache"
	"salon_be/component/ekycclient"
	"salon_be/component/sms"
	"salon_be/component/sms/esms"
	models "salon_be/model"
	pb "salon_be/proto/salon_be/salon_be"

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
	GetSMSClient() SMSClient
	GetEKYCClient() *ekycclient.EKYCClient
}

type DBInstances interface {
	GetMySQLDBConnection() *gorm.DB
	AutoMigrateMySQL() error
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
	RegisterServiceJobs(ctx context.Context, appCtx AppContext) error
}

type AppCache interface {
	GetUserCache(ctx context.Context, fakeUserId string) (string, error)
	SetUserCache(ctx context.Context, user *models.User) error
	DeleteUserCache(ctx context.Context, userId string) error

	SetVideoCache(ctx context.Context, serviceSlug string, video models.Video) error
	GetVideoCache(ctx context.Context, serviceSlug string, videoId string) (*cache.VideoCacheInfo, error)

	SetEnrollmentCache(ctx context.Context, enrollment *cache.EnrollmentCache) error
	GetEnrollmentCache(ctx context.Context, serviceSlug, userId string) (*cache.EnrollmentCache, error)
	DeleteEnrollmentCache(ctx context.Context, serviceSlug, userId string) error

	CreateDynamoDBTables(ctx context.Context, tables []cache.DynamoDBTableDefinition) error
	GetDynamoDBTableDefinitions() []cache.DynamoDBTableDefinition
}

type LocalPubSub interface {
	GetUnblockPubSub() *gochannel.GoChannel
	// The publisher of this PubSub will be block util its message got acked
	GetBlockPubSub() *gochannel.GoChannel
}

type SMSClient interface {
	SendOTP(ctx context.Context, otpMessage sms.OTPMessage) (*esms.ESMSResponse, error)
}

type appCtx struct {
	dbInstances        DBInstances
	localPubSub        LocalPubSub
	videoServiceClient pb.VideoProcessingServiceClient
	jwtSecretKey       string
	appCache           AppCache
	awsSession         *session.Session
	s3Client           *s3.S3
	appqueue           AppQueue
	cron               CronJob
	smsClient          SMSClient
	ekycClient         *ekycclient.EKYCClient // Added this
}

func NewAppContext(
	dbInstances DBInstances,
	localPubSub LocalPubSub,
	videoServiceClient pb.VideoProcessingServiceClient,
	jwtSecretKey string,
	appCache AppCache,
	awsSession *session.Session,
	appqueue AppQueue,
	cron CronJob,
	s3Client *s3.S3,
	smsClient SMSClient,
	ekycClient *ekycclient.EKYCClient,
) *appCtx {
	return &appCtx{
		dbInstances:        dbInstances,
		localPubSub:        localPubSub,
		videoServiceClient: videoServiceClient,
		jwtSecretKey:       jwtSecretKey,
		appCache:           appCache,
		awsSession:         awsSession,
		appqueue:           appqueue,
		cron:               cron,
		s3Client:           s3Client,
		smsClient:          smsClient,
		ekycClient:         ekycClient,
	}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.dbInstances.GetMySQLDBConnection()
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

func (ctx *appCtx) GetSMSClient() SMSClient {
	return ctx.smsClient
}

func (ctx *appCtx) GetEKYCClient() *ekycclient.EKYCClient {
	return ctx.ekycClient
}
