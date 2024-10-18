package grpcserver

import (
	"context"
	"fmt"
	"time"
	"video_server/component"
	"video_server/component/logger"
	pb "video_server/proto/video_service/video_service"
	"video_server/watermill/messagemodel"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToVideoProcessingServer(ctx context.Context) (pb.VideoProcessingServiceClient, *grpc.ClientConn, error) {
	videoProcessorAddr := "video-processor:50052"

	// Set up a connection to the server with a timeout
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, videoProcessorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to connect to video processing server", zap.Error(err))
		return nil, nil, fmt.Errorf("failed to connect to video processing server: %w", err)
	}

	// Create a new client
	videoServiceClient := pb.NewVideoProcessingServiceClient(conn)

	return videoServiceClient, conn, nil
}

func ProcessNewVideoRequest(ctx context.Context, appCtx component.AppContext, videoInfo *messagemodel.RequestProcessVideoInfo) error {
	// Prepare the request
	req := &pb.VideoInfo{
		S3Key:      videoInfo.RawVidS3Key,
		VideoId:    videoInfo.VideoId,
		CourseId:   videoInfo.CourseId,
		UploadedBy: videoInfo.UploadedBy,
	}

	// Create a timeout for the gRPC call
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Call the gRPC method
	resp, err := appCtx.GetVideoProcessingClient().ProcessNewVideoRequest(callCtx, req)
	if err != nil {
		logger.AppLogger.Error(ctx, "ProcessNewVideoRequest failed",
			zap.Any("req", req),
			zap.Error(err))
		return fmt.Errorf("ProcessNewVideoRequest failed: %w", err)
	}

	if resp.Status != codes.OK.String() {
		logger.AppLogger.Error(ctx, "ProcessNewVideoRequest returned non-OK status",
			zap.Any("req", req),
			zap.Any("resp", resp))
		return fmt.Errorf("ProcessNewVideoRequest returned non-OK status: %s", resp.Status)
	}

	// Handle the response
	logger.AppLogger.Info(ctx, "ProcessNewVideoRequest completed successfully",
		zap.Any("req", req),
		zap.Any("resp", resp))
	return nil
}
