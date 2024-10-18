package cronjob

import (
	"context"
	"video_server/component"
	"video_server/component/cronjob/job"
	"video_server/component/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getAllVideoJobs(ctx context.Context, appCtx component.AppContext) []Job {
	videoProcessReqInterval, err := MinutesToCronFormat(viper.GetInt("VIDEO_PROCESS_REQUEST_INTERVAL"))
	if err != nil {
		logger.AppLogger.Fatal(ctx, "failed calculating cronjob interval", zap.Error(err))
	}

	return []Job{
		{
			"RequestProcessVideoJob",
			videoProcessReqInterval,
			job.RequestProcessVideoJob(appCtx),
		},
	}
}

func (cj *cronJob) RegisterVideoJobs(ctx context.Context, appCtx component.AppContext) error {
	videoJobs := getAllVideoJobs(ctx, appCtx)
	for _, job := range videoJobs {
		_, err := cj.registerJob(ctx, cj.cron, job)
		if err != nil {
			logger.AppLogger.Fatal(ctx, "start cronjon failed", zap.Error(err))
		}
	}

	return nil
}
