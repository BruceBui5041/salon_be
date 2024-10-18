package lecturerepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lecture/lecturemodel"
)

type CreateLectureCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type CreateLectureStore interface {
	Create(
		ctx context.Context,
		newLecture *models.Lecture,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lecture, error)
}

type createLectureRepo struct {
	lectureStore CreateLectureStore
	courseStore  CreateLectureCourseStore
}

func NewCreateLectureRepo(lectureStore CreateLectureStore, courseStore CreateLectureCourseStore) *createLectureRepo {
	return &createLectureRepo{
		lectureStore: lectureStore,
		courseStore:  courseStore,
	}
}

func (repo *createLectureRepo) CreateNewLecture(ctx context.Context, input *lecturemodel.CreateLecture) (*models.Lecture, error) {
	courseUid, err := common.FromBase58(input.CourseID)
	if err != nil {
		return nil, err
	}

	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"id": courseUid.GetLocalID()})
	if err != nil {
		return nil, err
	}

	requester := ctx.Value(common.CurrentUser).(common.Requester)
	if course.CreatorID != requester.GetUserId() {
		return nil, common.ErrNoPermission(nil)
	}

	newLecture := &models.Lecture{
		CourseID:    course.Id,
		Title:       input.Title,
		Description: input.Description,
		Order:       input.Order,
	}

	lectureId, err := repo.lectureStore.Create(ctx, newLecture)
	if err != nil {
		return nil, err
	}

	lecture, err := repo.lectureStore.FindOne(ctx, map[string]interface{}{"id": lectureId})
	if err != nil {
		return nil, err
	}

	return lecture, nil
}
