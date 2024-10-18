// File: lecturerepo/updatelecture.go

package lecturerepo

import (
	"context"
	models "video_server/model"
	"video_server/model/lecture/lecturemodel"
)

type UpdateLectureStore interface {
	Update(
		ctx context.Context,
		lectureId uint32,
		data *lecturemodel.UpdateLecture,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lecture, error)
}

type UpdateLectureCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type updateLectureRepo struct {
	lectureStore UpdateLectureStore
	courseStore  UpdateLectureCourseStore
}

func NewUpdateLectureRepo(lectureStore UpdateLectureStore, courseStore UpdateLectureCourseStore) *updateLectureRepo {
	return &updateLectureRepo{
		lectureStore: lectureStore,
		courseStore:  courseStore,
	}
}

func (repo *updateLectureRepo) UpdateLecture(ctx context.Context, lectureId uint32, input *lecturemodel.UpdateLecture) error {
	return repo.lectureStore.Update(ctx, lectureId, input)
}

func (repo *updateLectureRepo) FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error) {
	return repo.courseStore.FindOne(ctx, conditions)
}

func (repo *updateLectureRepo) FindLecture(ctx context.Context, conditions map[string]interface{}) (*models.Lecture, error) {
	return repo.lectureStore.FindOne(ctx, conditions)
}
