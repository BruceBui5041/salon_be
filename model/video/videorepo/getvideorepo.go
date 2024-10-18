package videorepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
)

type GetVideoCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type GetVideoStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Video, error)
}

type getVideoRepo struct {
	videoStore  GetVideoStore
	courseStore GetVideoCourseStore
}

func NewGetVideoRepo(videoStore GetVideoStore, courseStore GetVideoCourseStore) *getVideoRepo {
	return &getVideoRepo{
		videoStore:  videoStore,
		courseStore: courseStore,
	}
}

func (repo *getVideoRepo) GetVideo(ctx context.Context, id uint32, courseSlug string) (*models.Video, error) {
	video, err := repo.videoStore.FindOne(ctx, map[string]interface{}{"id": id}, "Course", "Lesson")
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.VideoEntityName, err)
	}

	if video.Course.Slug != courseSlug {
		return nil, errors.New("course slug mismatch")
	}

	return video, nil
}
