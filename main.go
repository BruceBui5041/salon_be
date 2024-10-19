package main

import (
	"context"
	"log"
	"os"
	"salon_be/component"
	"salon_be/component/config"
	"salon_be/component/db"
	"salon_be/component/logger"
	"salon_be/component/server"
	"salon_be/component/telemetry"
	"salon_be/watermill"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	shutdown := telemetry.InitTracer()
	defer shutdown()

	tracer := otel.Tracer("video_service_start")
	ctx, span := tracer.Start(ctx, "video_service_start", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	logger.CreateAppLogger(ctx)

	jwtSecretKey := os.Getenv("JWTSecretKey")

	// client, conn, err := grpcserver.ConnectToVideoProcessingServer(ctx)
	// if err != nil {
	// 	logger.AppLogger.Fatal(ctx, "failed to start grpc server", zap.Error(err))
	// }
	// defer conn.Close()

	awsSession, err := config.CreateAWSSession()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	appContext := component.NewAppContext(
		db.ConnectToDB(ctx),
		watermill.NewPubsubPublisher(),
		nil,
		jwtSecretKey,
		config.CreateAppCache(awsSession),
		awsSession,
		config.CreateAppQueue(awsSession),
		// cronjob.CreateCron(),
		nil,
		config.CreateS3Client(awsSession),
	)

	// appCache := appContext.GetAppCache()
	// if err := appCache.CreateDynamoDBTables(ctx, appCache.GetDynamoDBTableDefinitions()); err != nil {
	// 	logger.AppLogger.Fatal(ctx, "create dynamaDB tables failed", zap.Error(err))
	// }

	// results := appContext.GetAppQueue().CreateSQSQueues(ctx)
	// logger.AppLogger.Info(ctx, "created sqs queue", zap.Any("res", results))

	go watermill.StartSubscribers(appContext)

	// Start gRPC server
	// startGRPCServer()

	// videoProcessProgressTopic := viper.GetString("UPDATE_VIDEO_PROCESS_PROGRESS_TOPIC")
	// consumeTopics := []string{videoProcessProgressTopic}
	// go appContext.GetAppQueue().StartSQSMessageListener(ctx, appContext, consumeTopics, consumerhandler.QueueMsgHander)

	// startCronJobs(appContext)

	// Start HTTP server
	server.StartHTTPServer(appContext)
}
