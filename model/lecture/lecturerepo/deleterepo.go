package lecturerepo

import (
	"context"
	models "video_server/model"
)

type DeleteLectureStore interface {
	Delete(ctx context.Context, lectureId uint32) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lecture, error)
}

type DeleteLectureCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type deleteLectureRepo struct {
	lectureStore DeleteLectureStore
	courseStore  DeleteLectureCourseStore
}

func NewDeleteLectureRepo(lectureStore DeleteLectureStore, courseStore DeleteLectureCourseStore) *deleteLectureRepo {
	return &deleteLectureRepo{
		lectureStore: lectureStore,
		courseStore:  courseStore,
	}
}

func (repo *deleteLectureRepo) DeleteLecture(ctx context.Context, lectureId uint32) error {
	return repo.lectureStore.Delete(ctx, lectureId)
}

func (repo *deleteLectureRepo) FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error) {
	return repo.courseStore.FindOne(ctx, conditions)
}

func (repo *deleteLectureRepo) FindLecture(ctx context.Context, conditions map[string]interface{}) (*models.Lecture, error) {
	return repo.lectureStore.FindOne(ctx, conditions)
}
