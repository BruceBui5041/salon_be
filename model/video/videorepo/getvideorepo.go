package videorepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
)

type GetVideoServiceStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Service, error)
}

type GetVideoStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Video, error)
}

type getVideoRepo struct {
	videoStore   GetVideoStore
	serviceStore GetVideoServiceStore
}

func NewGetVideoRepo(videoStore GetVideoStore, serviceStore GetVideoServiceStore) *getVideoRepo {
	return &getVideoRepo{
		videoStore:   videoStore,
		serviceStore: serviceStore,
	}
}

func (repo *getVideoRepo) GetVideo(ctx context.Context, id uint32, serviceSlug string) (*models.Video, error) {
	video, err := repo.videoStore.FindOne(ctx, map[string]interface{}{"id": id}, "Service", "Lesson")
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.VideoEntityName, err)
	}

	if video.Service.Slug != serviceSlug {
		return nil, errors.New("service slug mismatch")
	}

	return video, nil
}
