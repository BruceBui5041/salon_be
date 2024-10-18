package progressbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/progress/progressmodel"
)

type ProgressRepo interface {
	CreateNewProgress(ctx context.Context, input *progressmodel.CreateProgress) (*models.Progress, error)
}

type createProgressBiz struct {
	repo ProgressRepo
}

func NewCreateProgressBiz(repo ProgressRepo) *createProgressBiz {
	return &createProgressBiz{repo: repo}
}

func (c *createProgressBiz) CreateNewProgress(ctx context.Context, input *progressmodel.CreateProgress) error {
	if input.UserID == 0 {
		return common.ErrInvalidRequest(errors.New("user ID is required"))
	}

	if input.VideoID == "" {
		return common.ErrInvalidRequest(errors.New("video ID is required"))
	}

	progress, err := c.repo.CreateNewProgress(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.ProgressEntityName, err)
	}

	input.Id = progress.Id

	return nil
}
