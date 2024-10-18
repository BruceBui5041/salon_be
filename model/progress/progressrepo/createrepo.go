package progressrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/progress/progressmodel"
)

type CreateProgressStore interface {
	Create(
		ctx context.Context,
		newProgress *models.Progress,
	) (*models.Progress, error)
}

type createProgressRepo struct {
	store CreateProgressStore
}

func NewCreateProgressRepo(store CreateProgressStore) *createProgressRepo {
	return &createProgressRepo{
		store: store,
	}
}

func (repo *createProgressRepo) CreateNewProgress(
	ctx context.Context,
	input *progressmodel.CreateProgress,
) (*models.Progress, error) {
	videoUID, err := common.FromBase58(input.VideoID)
	if err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	newProgress := &models.Progress{
		UserID:         input.UserID,
		VideoID:        videoUID.GetLocalID(),
		WatchedSeconds: input.WatchedSeconds,
		Completed:      input.Completed,
		LastWatched:    input.LastWatched,
	}

	progress, err := repo.store.Create(ctx, newProgress)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return progress, nil
}
