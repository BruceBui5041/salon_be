package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"video_server/apihandler"
	"video_server/common"
	"video_server/component"
	"video_server/component/appqueue"
	"video_server/component/appqueue/consumerhandler"
	"video_server/component/cache"
	"video_server/component/cronjob"
	"video_server/component/genericapi/generictransport"
	"video_server/component/grpcserver"
	"video_server/component/logger"
	"video_server/component/telemetry"
	"video_server/component/ws"
	"video_server/middleware"
	"video_server/model/category/categorytransport"
	"video_server/model/comment/commenttransport"
	"video_server/model/course/coursetransport"
	"video_server/model/lecture/lecturetransport"
	"video_server/model/lesson/lessontransport"
	"video_server/model/payment/paymenttransport"
	"video_server/model/permission/permissiontransport"
	"video_server/model/role/roletransport"
	"video_server/model/user/usertransport"
	"video_server/model/userprofile/userprofiletransport"
	"video_server/model/video/videotransport"
	"video_server/watermill"

	pb "video_server/proto/video_service/video_service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

	awsSession, err := createAWSSession()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	appContext := component.NewAppContext(
		connectToDB(ctx),
		watermill.NewPubsubPublisher(),
		nil,
		jwtSecretKey,
		createAppCache(awsSession),
		awsSession,
		appqueue.CreateAppQueue(awsSession),
		cronjob.CreateCron(),
		s3.New(awsSession),
	)

	appCache := appContext.GetAppCache()
	if err := appCache.CreateDynamoDBTables(ctx, appCache.GetDynamoDBTableDefinitions()); err != nil {
		logger.AppLogger.Fatal(ctx, "create dynamaDB tables failed", zap.Error(err))
	}

	results := appContext.GetAppQueue().CreateSQSQueues(ctx)
	logger.AppLogger.Info(ctx, "created sqs queue", zap.Any("res", results))

	go watermill.StartSubscribers(appContext)

	// Start gRPC server
	// startGRPCServer()

	videoProcessProgressTopic := viper.GetString("UPDATE_VIDEO_PROCESS_PROGRESS_TOPIC")
	consumeTopics := []string{videoProcessProgressTopic}
	go appContext.GetAppQueue().StartSQSMessageListener(ctx, appContext, consumeTopics, consumerhandler.QueueMsgHander)

	startCronJobs(appContext)

	// Start HTTP server
	startHTTPServer(appContext)

}

func connectToDB(ctx context.Context) *gorm.DB {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormlogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		logger.AppLogger.Fatal(ctx, err.Error())
	}

	return db
}

func startHTTPServer(appCtx component.AppContext) {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:8083"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.MaxAge = 300

	r.Use(cors.New(config))
	r.Use(middleware.OtelTracing())
	r.Use(middleware.Recover(appCtx))

	wsServer := ws.NewWebSocketServer()

	r.GET("/ws", wsServer.HandleWebSocket(appCtx))

	genTransport := generictransport.NewGenericTransport(appCtx)
	genAPIs := r.Group("/")
	{
		genAPIs.POST("search", genTransport.Search())
		// genAPIs.POST("create", middleware.RequiredAuth(appCtx), genTransport.Create())
	}

	permission := r.Group("/permission", middleware.RequiredAuth(appCtx))
	{
		permission.POST("", permissiontransport.CreatePermissionHandler(appCtx))
		permission.PATCH("/:id", permissiontransport.UpdatePermissionHandler(appCtx))
	}

	role := r.Group("/role", middleware.RequiredAuth(appCtx))
	{
		role.POST("", roletransport.CreateRoleHandler(appCtx))
		role.PATCH("/:id", roletransport.UpdateRoleHandler(appCtx))
		role.DELETE("/:id", roletransport.SoftDeleteRoleHandler(appCtx))
	}

	user := r.Group("/user")
	{
		user.GET("", middleware.RequiredAuth(appCtx), usertransport.GetUser(appCtx))
		user.PATCH("/:id", middleware.RequiredAuth(appCtx), usertransport.UpdateUser(appCtx))
	}

	userprofile := r.Group("/profile")
	{
		userprofile.POST("", middleware.RequiredAuth(appCtx), userprofiletransport.CreateUserProfileHandler(appCtx))
		userprofile.PUT("", middleware.RequiredAuth(appCtx), userprofiletransport.UpdateProfileHandler(appCtx))
	}

	previewVideo := r.Group("/preview_video", middleware.PublicVideoCheck(appCtx))
	{
		// for get master list
		previewVideo.GET("/playlist/:course_slug/:video_id",
			apihandler.GetPlaylistHandler(appCtx),
		)

		// for get video playlish
		previewVideo.GET(
			"/playlist/:course_slug/:video_id/:resolution/:playlistName",
			apihandler.GetPlaylistHandler(appCtx),
		)

		previewVideo.GET("",
			apihandler.SegmentHandler(appCtx),
		)
	}

	commentGroup := r.Group("/comment")
	{
		commentGroup.POST("", middleware.RequiredAuth(appCtx), commenttransport.CreateCommentHandler(appCtx))
		commentGroup.PUT("/:id", middleware.RequiredAuth(appCtx), commenttransport.UpdateCommentHandler(appCtx))
	}

	courseGroup := r.Group("/course")
	{
		courseGroup.GET("", coursetransport.ListCourses(appCtx))
		courseGroup.GET("/:id", coursetransport.GetCourseByID(appCtx))

	}

	courseGroupInstructor := r.Group(
		"/course",
		middleware.RequiredAuth(appCtx),
		middleware.AllowIntructorOnly(appCtx),
	)
	{
		courseGroupInstructor.POST("", coursetransport.CreateCourseHandler(appCtx))
		courseGroupInstructor.PUT("/:id", coursetransport.UpdateCourseHandler(appCtx))
		courseGroupInstructor.GET("/:id/videos", coursetransport.GetCourseVideos(appCtx))
	}

	lessonGroup := r.Group("/lessons", middleware.RequiredAuth(appCtx), middleware.AllowIntructorOnly(appCtx))
	{
		lessonGroup.POST("", lessontransport.CreateLessonHandler(appCtx))
		lessonGroup.PUT("/:id", lessontransport.UpdateLessonHandler(appCtx))
		lessonGroup.DELETE("/:id", lessontransport.DeleteLessonHandler(appCtx))
		lessonGroup.GET("/:id", lessontransport.GetLessonHandler(appCtx))
	}

	lectureGroup := r.Group(
		"/lectures",
		middleware.RequiredAuth(appCtx),
		middleware.AllowIntructorOnly(appCtx),
	)
	{
		lectureGroup.POST("", lecturetransport.CreateLectureHandler(appCtx))
		lectureGroup.PUT("/:id", lecturetransport.UpdateLectureHandler(appCtx))
		lectureGroup.DELETE("/:id", lecturetransport.DeleteLectureHandler(appCtx))
	}

	videoGroupInstructor := r.Group("/video",
		middleware.RequiredAuth(appCtx),
		middleware.AllowIntructorOnly(appCtx),
	)
	{
		videoGroupInstructor.POST("", videotransport.CreateVideoHandler(appCtx))
		videoGroupInstructor.PUT("/:id", videotransport.UpdateVideoHandler(appCtx))
	}

	video := r.Group(
		"/video",
		middleware.RequiredAuth(appCtx),
		middleware.EnrollmentCheck(appCtx),
	)
	{
		// for get master list
		video.GET("/playlist/:course_slug/:video_id", apihandler.GetPlaylistHandler(appCtx))

		// for get video playlish
		video.GET(
			"/playlist/:course_slug/:video_id/:resolution/:playlistName",
			apihandler.GetPlaylistHandler(appCtx),
		)

		video.GET("", apihandler.SegmentHandler(appCtx))
		// video.GET("/:id", videotransport.GetVideoBySlug(appCtx))
		videoGroupInstructor.GET(
			"/:course_slug",
			videotransport.ListCourseVideos(appCtx),
		)

	}

	paymentGroup := r.Group("/payment", middleware.RequiredAuth(appCtx))
	{
		paymentGroup.POST("", paymenttransport.CreatePaymentHandler(appCtx))
	}

	categoryGroup := r.Group("/category")
	{
		categoryGroup.GET("", categorytransport.ListCategories(appCtx))
		categoryGroup.PATCH("/:id", middleware.RequiredAuth(appCtx), categorytransport.UpdateCategoryHandler(appCtx))
		categoryGroup.POST("", middleware.RequiredAuth(appCtx), categorytransport.CreateCategoryHandler(appCtx))
	}

	r.GET("/checkauth", middleware.RequiredAuth(appCtx), func(c *gin.Context) {
		requester, ok := c.Request.Context().Value(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}
		var userCached cache.CacheUser
		copier.Copy(&userCached, requester)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(userCached))
	})

	r.POST("/login", usertransport.Login(appCtx))
	r.POST("/register", usertransport.Register(appCtx))
	r.POST("/logout", middleware.RequiredAuth(appCtx), usertransport.Logout(appCtx))

	// TODO: disable this apis if not dev env
	r.GET("/decode/:id", apihandler.DecodeUID(appCtx))
	r.GET("/encode/:id/:dbtype", apihandler.EncodeUID(appCtx))

	// r.GET("test", videotransport.CreateVideoHandlerTest(appCtx))

	restPort := viper.GetString("REST_PORT")
	log.Printf("Starting HTTP server on :%s", restPort)
	if err := r.Run(fmt.Sprintf(":%s", restPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register your gRPC services here
	pb.RegisterVideoProcessingServiceServer(s, &grpcserver.VideoServiceServer{})

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func createAWSSession() (*session.Session, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return sess, nil
}

func createAppCache(awsSess *session.Session) component.AppCache {
	appcache, err := cache.NewAppCache(awsSess)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	return appcache
}

func startCronJobs(appCtx component.AppContext) {
	cronCtx := context.Background()
	tracer := otel.Tracer("CRONJOB")
	cronCtx, span := tracer.Start(cronCtx, "cron job update course count field", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	appCron := appCtx.GetCronJob()
	appCron.RegisterVideoJobs(cronCtx, appCtx)
	appCron.RegisterCourseJobs(cronCtx, appCtx)
	appCron.Start()
}
