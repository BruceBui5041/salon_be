package cronjob

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Job struct {
	jobName string
	spec    string
	cmd     func()
}

type cronJob struct {
	cron              *cron.Cron
	mapJobIdToJobName map[cron.EntryID]string
}

func CreateCron() *cronJob {
	c := &cronJob{
		cron:              cron.New(),
		mapJobIdToJobName: map[cron.EntryID]string{},
	}
	return c
}

func (cj *cronJob) Start() {
	cj.cron.Start()
}

func (cj *cronJob) registerJob(ctx context.Context, videoCron *cron.Cron, job Job) (cron.EntryID, error) {
	entryId, err := videoCron.AddFunc(job.spec, WithDurationLogging(ctx, job.jobName, job.cmd))
	cj.mapJobIdToJobName[entryId] = job.jobName
	return entryId, err
}
