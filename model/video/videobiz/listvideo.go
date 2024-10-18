package videobiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type ListVideoRepo interface {
	ListCourseVideos(ctx context.Context, courseSlug string) ([]models.Video, error)
}

type listVideoBiz struct {
	listVideoRepo ListVideoRepo
}

func NewListVideoBiz(repo ListVideoRepo) *listVideoBiz {
	return &listVideoBiz{listVideoRepo: repo}
}

func (biz *listVideoBiz) ListCourseVideos(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Video, error) {
	courseSlug := conditions["course_slug"].(string)
	videos, err := biz.listVideoRepo.ListCourseVideos(ctx, courseSlug)
	if err != nil {
		return nil, common.ErrCannotListEntity(models.VideoEntityName, err)
	}
	return videos, nil
}
