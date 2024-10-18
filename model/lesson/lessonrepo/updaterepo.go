// lesson/lessonrepo/updaterepo.go

package lessonrepo

import (
	"context"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"
	"video_server/model/video/videomodel"
)

type UpdateLessonStore interface {
	UpdateLesson(
		ctx context.Context,
		lessonId uint32,
		data *lessonmodel.UpdateLesson,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lesson, error)
}

type UpdateLessonCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type UpdateVideoStore interface {
	UpdateVideo(
		ctx context.Context,
		id uint32,
		updateData *videomodel.UpdateVideo,
	) error
}

type updateLessonRepo struct {
	lessonStore UpdateLessonStore
	courseStore UpdateLessonCourseStore
	videoStore  UpdateVideoStore
}

func NewUpdateLessonRepo(
	lessonStore UpdateLessonStore,
	courseStore UpdateLessonCourseStore,
	videoStore UpdateVideoStore,
) *updateLessonRepo {
	return &updateLessonRepo{
		lessonStore: lessonStore,
		courseStore: courseStore,
		videoStore:  videoStore,
	}
}

func (repo *updateLessonRepo) UpdateLesson(ctx context.Context, lessonId uint32, input *lessonmodel.UpdateLesson) error {
	return repo.lessonStore.UpdateLesson(ctx, lessonId, input)
}

func (repo *updateLessonRepo) FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error) {
	return repo.courseStore.FindOne(ctx, conditions)
}

func (repo *updateLessonRepo) FindLesson(ctx context.Context, conditions map[string]interface{}) (*models.Lesson, error) {
	return repo.lessonStore.FindOne(ctx, conditions)
}

func (repo *updateLessonRepo) UpdateVideo(ctx context.Context, videoId uint32, input *videomodel.UpdateVideo) error {
	return repo.videoStore.UpdateVideo(ctx, videoId, input)
}
