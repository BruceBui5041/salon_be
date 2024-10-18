package lessonrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"
)

type CreateLessonCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type CreateLessonVideoStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Video, error)
}

type CreateLessonStore interface {
	CreateNewLesson(
		ctx context.Context,
		newLesson *models.Lesson,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lesson, error)
}

type createLessonRepo struct {
	lessonStore CreateLessonStore
	courseStore CreateLessonCourseStore
	videoStore  CreateLessonVideoStore
}

func NewCreateLessonRepo(lessonStore CreateLessonStore, courseStore CreateLessonCourseStore, videoStore CreateLessonVideoStore) *createLessonRepo {
	return &createLessonRepo{
		lessonStore: lessonStore,
		courseStore: courseStore,
		videoStore:  videoStore,
	}
}

func (repo *createLessonRepo) CreateNewLesson(ctx context.Context, input *lessonmodel.CreateLesson) (*models.Lesson, error) {
	courseUid, err := common.FromBase58(input.CourseID)
	if err != nil {
		return nil, err
	}

	lectureUid, err := common.FromBase58(input.LectureID)
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

	newLesson := &models.Lesson{
		CourseID:     course.Id,
		LectureID:    lectureUid.GetLocalID(),
		Title:        input.Title,
		Description:  input.Description,
		Duration:     input.Duration,
		Order:        input.Order,
		Type:         input.Type,
		AllowPreview: input.AllowPreview,
	}

	lessonId, err := repo.lessonStore.CreateNewLesson(ctx, newLesson)
	if err != nil {
		return nil, err
	}

	lesson, err := repo.lessonStore.FindOne(ctx, map[string]interface{}{"id": lessonId})
	if err != nil {
		return nil, err
	}

	return lesson, nil
}
