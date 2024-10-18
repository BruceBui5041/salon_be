package cronjob

import (
	"context"
	"video_server/component"
	"video_server/component/cronjob/job"
	"video_server/component/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getAllCourseJobs(ctx context.Context, appCtx component.AppContext) []Job {
	courseCountStudentInterval, err := MinutesToCronFormat(viper.GetInt("COURSE_UPDATE_COUNTFIELD_INTERVAL"))
	if err != nil {
		logger.AppLogger.Fatal(ctx, "failed calculating cronjob interval", zap.Error(err))
	}

	return []Job{
		{
			"CourseUpdateCountFieldJob",
			courseCountStudentInterval,
			job.CourseUpdateCountFieldJob(ctx, appCtx),
		},
	}
}

func (cj *cronJob) RegisterCourseJobs(ctx context.Context, appCtx component.AppContext) error {
	courseJobs := getAllCourseJobs(ctx, appCtx)
	for _, job := range courseJobs {
		_, err := cj.registerJob(ctx, cj.cron, job)
		if err != nil {
			logger.AppLogger.Fatal(ctx, "start cronjob failed", zap.Error(err))
		}
	}

	return nil
}
