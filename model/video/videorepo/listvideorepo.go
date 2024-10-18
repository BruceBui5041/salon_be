package videorepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type VideoStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Video, error)
}

type ListCourseStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...interface{}) ([]models.Course, error)
}

type listVideoRepo struct {
	videoStore  VideoStore
	courseStore ListCourseStore
}

func NewListVideoRepo(videoStore VideoStore, listCourseStore ListCourseStore) *listVideoRepo {
	return &listVideoRepo{
		videoStore:  videoStore,
		courseStore: listCourseStore,
	}
}

func (repo *listVideoRepo) ListCourseVideos(ctx context.Context, courseSlug string) ([]models.Video, error) {
	courseConditions := map[string]interface{}{"slug": courseSlug}
	courses, err := repo.courseStore.FindAll(ctx, courseConditions)
	if err != nil {
		return nil, err
	}

	if len(courses) == 0 {
		return nil, common.RecordNotFound
	}

	videoConditions := map[string]interface{}{"course_id": courses[0].Id}
	videos, err := repo.videoStore.Find(ctx, videoConditions, "ProcessInfos")
	if err != nil {
		return nil, err
	}

	return videos, nil
}
