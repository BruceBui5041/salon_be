package videorepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type VideoStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Video, error)
}

type ListServiceStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...interface{}) ([]models.ServiceVersion, error)
}

type listVideoRepo struct {
	videoStore   VideoStore
	serviceStore ListServiceStore
}

func NewListVideoRepo(videoStore VideoStore, listServiceStore ListServiceStore) *listVideoRepo {
	return &listVideoRepo{
		videoStore:   videoStore,
		serviceStore: listServiceStore,
	}
}

func (repo *listVideoRepo) ListServiceVideos(ctx context.Context, serviceSlug string) ([]models.Video, error) {
	serviceConditions := map[string]interface{}{"slug": serviceSlug}
	services, err := repo.serviceStore.FindAll(ctx, serviceConditions)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, common.RecordNotFound
	}

	videoConditions := map[string]interface{}{"service_id": services[0].Id}
	videos, err := repo.videoStore.Find(ctx, videoConditions, "ProcessInfos")
	if err != nil {
		return nil, err
	}

	return videos, nil
}
