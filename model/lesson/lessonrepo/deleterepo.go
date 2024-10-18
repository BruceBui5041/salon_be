package lessonrepo

import (
	"context"
	models "video_server/model"
)

type DeleteLessonStore interface {
	Delete(ctx context.Context, lessonId uint32) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lesson, error)
}

type DeleteLessonCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type deleteLessonRepo struct {
	lessonStore DeleteLessonStore
	courseStore DeleteLessonCourseStore
}

func NewDeleteLessonRepo(lessonStore DeleteLessonStore, courseStore DeleteLessonCourseStore) *deleteLessonRepo {
	return &deleteLessonRepo{
		lessonStore: lessonStore,
		courseStore: courseStore,
	}
}

func (repo *deleteLessonRepo) DeleteLesson(ctx context.Context, lessonId uint32) error {
	return repo.lessonStore.Delete(ctx, lessonId)
}

func (repo *deleteLessonRepo) FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error) {
	return repo.courseStore.FindOne(ctx, conditions)
}

func (repo *deleteLessonRepo) FindLesson(ctx context.Context, conditions map[string]interface{}) (*models.Lesson, error) {
	return repo.lessonStore.FindOne(ctx, conditions)
}
