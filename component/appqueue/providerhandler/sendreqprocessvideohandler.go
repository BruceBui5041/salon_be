package providerhandler

import (
	"context"
	"encoding/base64"
	"video_server/appconst"
	"video_server/component"
	"video_server/component/logger"
	pb "video_server/proto/video_service/video_service"
	"video_server/watermill/messagemodel"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func SendRequestProcessVideo(ctx context.Context, appQ component.AppQueue, videoInfo *messagemodel.RequestProcessVideoInfo) error {
	req := &pb.VideoInfo{
		S3Key:             videoInfo.RawVidS3Key,
		VideoId:           videoInfo.VideoId,
		CourseId:          videoInfo.CourseId,
		UploadedBy:        videoInfo.UploadedBy,
		Timestamp:         videoInfo.Timestamp,
		Retry:             int32(videoInfo.Retry),
		RequestResolution: *videoInfo.RequestResolution.Enum(),
	}

	logger.AppLogger.Info(ctx, "RequestProcessVideoInfo", zap.Any("req info", videoInfo))
	queueMsg, err := proto.Marshal(req)
	if err != nil {
		logger.AppLogger.Error(ctx, "Cannot unmarshal video info protobuf", zap.Any("protobuf", req), zap.Error(err))
		return err
	}

	// Encode the serialized data to Base64
	base64Encoded := base64.StdEncoding.EncodeToString(queueMsg)

	reqProcessVideoTopic := viper.GetString("REQ_PROCESS_VIDEO_TOPIC")
	err = appQ.SendSQSMessage(
		ctx,
		reqProcessVideoTopic,
		appconst.ReqProcessVideoGroupId,
		base64Encoded,
	)

	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"Cannot send request process video message",
			zap.Any("topic", reqProcessVideoTopic),
			zap.Any("msg base64", base64Encoded),
			zap.Error(err),
		)
		return err
	}

	return nil
}
