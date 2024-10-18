package videobiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type ListVideoRepo interface {
	ListServiceVideos(ctx context.Context, serviceSlug string) ([]models.Video, error)
}

type listVideoBiz struct {
	listVideoRepo ListVideoRepo
}

func NewListVideoBiz(repo ListVideoRepo) *listVideoBiz {
	return &listVideoBiz{listVideoRepo: repo}
}

func (biz *listVideoBiz) ListServiceVideos(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Video, error) {
	serviceSlug := conditions["service_slug"].(string)
	videos, err := biz.listVideoRepo.ListServiceVideos(ctx, serviceSlug)
	if err != nil {
		return nil, common.ErrCannotListEntity(models.VideoEntityName, err)
	}
	return videos, nil
}
