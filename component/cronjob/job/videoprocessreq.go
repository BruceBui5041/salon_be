package job

import (
	"context"
	"fmt"
	"time"
	"video_server/common"
	"video_server/component"
	"video_server/component/appqueue/providerhandler"
	"video_server/component/logger"
	"video_server/model/videoprocessinfo/videoprocessinfostore"
	"video_server/utils"
	"video_server/watermill/messagemodel"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RequestProcessVideoJob(appCtx component.AppContext) func() {
	return func() {
		// Create context when the job is actually executed
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		tracer := otel.Tracer("CRONJOB")
		ctx, span := tracer.Start(ctx, "cron job update course count field", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			videoProcessStore := videoprocessinfostore.NewSQLStore(tx)
			processInfo, err := videoProcessStore.Find(
				ctx,
				nil,
				"process_state IN ('pending', 'error') AND process_retry < 6",
				"Video.Course",
			)

			if err != nil {
				logger.AppLogger.Error(ctx, "RequestProcessVideoJob failed to find processinfo", zap.Error(err))
				return err
			}

			for _, pInfo := range processInfo {
				utcTime := time.Now().UTC()
				timestamp := fmt.Sprintf("%d", utcTime.UnixNano())

				pInfo.Video.Mask(false)
				pInfo.Video.Course.Mask(false)

				sqlModel := common.SQLModel{Id: pInfo.Video.Course.CreatorID}
				sqlModel.GenUID(common.DbTypeUser)

				protoResolution, err := utils.StringToProcessResolution(pInfo.ProcessResolution)
				if err != nil {
					logger.AppLogger.Error(
						ctx,
						"RequestProcessVideoJob failed to parse string to protobuf resolution",
						zap.Error(err),
						zap.Any("pInfo.ProcessResolution", pInfo.ProcessResolution),
					)
					return err
				}

				videoInfo := &messagemodel.RequestProcessVideoInfo{
					Timestamp:         timestamp,
					RawVidS3Key:       pInfo.Video.RawVideoURL,
					UploadedBy:        sqlModel.GetFakeId(),
					CourseId:          pInfo.Video.Course.GetFakeId(),
					VideoId:           pInfo.Video.GetFakeId(),
					Retry:             uint(pInfo.ProcessRetry),
					RequestResolution: &protoResolution,
				}

				err = providerhandler.SendRequestProcessVideo(ctx, appCtx.GetAppQueue(), videoInfo)

				if err != nil {
					logger.AppLogger.Error(
						ctx,
						"RequestProcessVideoJob failed to send sqs request process video",
						zap.Error(err),
						zap.Any("videoInfo", videoInfo),
					)
					return err
				}

				pInfo.ProcessState = "inqueue"
				if pInfo.ProcessRetry > 0 {
					pInfo.ProcessRetry = pInfo.ProcessRetry + 1
				}
			}

			err = videoProcessStore.UpdateMultiProcessState(ctx, processInfo)
			if err != nil {
				logger.AppLogger.Error(
					ctx,
					"fail to update process state",
					zap.Error(err),
					zap.Any("update processInfo", processInfo),
				)
				return err
			}
			return nil
		}); err != nil {
			logger.AppLogger.Error(ctx, "RequestProcessVideoJob failed to process processinfo", zap.Error(err))
		}
	}
}
