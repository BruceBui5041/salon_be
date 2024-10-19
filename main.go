package main

import (
	"context"
	"flag"
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
		nil,
		config.CreateS3Client(awsSession),
	)

	go watermill.StartSubscribers(appContext)

	// Start HTTP server
	server.StartHTTPServer(appContext)
}
