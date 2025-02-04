package main

import (
	"context"
	"flag"
	"log"
	"os"
	"salon_be/component"
	"salon_be/component/config"
	"salon_be/component/db"
	"salon_be/component/ekycclient"
	"salon_be/component/logger"
	"salon_be/component/server"
	"salon_be/component/sms"
	"salon_be/component/telemetry"
	"salon_be/watermill"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	// Parse command-line flags
	updateDB := flag.String("udb", "", "Update database schema (use 'mysql' to update MySQL)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// Database migration mode
	if *updateDB == "mysql" {
		log.Println("Connecting to MySQL database...")
		dbInstances := db.ConnectToDB(ctx)

		log.Println("Checking and updating MySQL database schema...")
		if err := dbInstances.AutoMigrateMySQL(); err != nil {
			log.Fatalf("Failed to check and update MySQL schema: %v", err)
		}
		log.Println("MySQL database schema check and update completed successfully")
		return // Exit after updating the database
	}

	// Normal service mode
	shutdown := telemetry.InitTracer()
	defer shutdown()

	tracer := otel.Tracer("salon_be_start")
	ctx, span := tracer.Start(ctx, "salon_be_start", trace.WithSpanKind(trace.SpanKindServer))
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
		sms.NewSMSClient(),
		ekycclient.NewEKYCClient(),
	)

	appCache := appContext.GetAppCache()
	if err := appCache.CreateDynamoDBTables(ctx, appCache.GetDynamoDBTableDefinitions()); err != nil {
		logger.AppLogger.Fatal(ctx, "create dynamaDB tables failed", zap.Error(err))
	}

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
