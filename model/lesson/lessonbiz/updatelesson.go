package lessonbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"
	"video_server/model/video/videomodel"
)

type UpdateLessonRepo interface {
	UpdateLesson(ctx context.Context, lessonId uint32, input *lessonmodel.UpdateLesson) error
	FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error)
	FindLesson(ctx context.Context, conditions map[string]interface{}) (*models.Lesson, error)
	UpdateVideo(ctx context.Context, videoId uint32, input *videomodel.UpdateVideo) error
}

type updateLessonBiz struct {
	repo UpdateLessonRepo
}

func NewUpdateLessonBiz(repo UpdateLessonRepo) *updateLessonBiz {
	return &updateLessonBiz{repo: repo}
}

func (biz *updateLessonBiz) UpdateLesson(ctx context.Context, lessonId string, input *lessonmodel.UpdateLesson) error {
	lessonUid, err := common.FromBase58(lessonId)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	lesson, err := biz.repo.FindLesson(ctx, map[string]interface{}{"id": lessonUid.GetLocalID()})
	if err != nil {
		return common.ErrCannotGetEntity(models.LessonEntityName, err)
	}

	if lesson == nil {
		return common.ErrEntityNotFound(models.LessonEntityName, errors.New("lesson not found"))
	}

	course, err := biz.repo.FindCourse(ctx, map[string]interface{}{"id": lesson.CourseID})
	if err != nil {
		return common.ErrCannotGetEntity(models.CourseEntityName, err)
	}

	requester, ok := ctx.Value(common.CurrentUser).(common.Requester)
	if !ok {
		return common.ErrNoPermission(errors.New("user not found"))
	}

	if course.CreatorID != requester.GetUserId() {
		return common.ErrNoPermission(errors.New("you don't have permission to update this lesson"))
	}

	if err := biz.repo.UpdateLesson(ctx, lesson.Id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.LessonEntityName, err)
	}

	if input.VideoID != nil {
		videoUpdate := &videomodel.UpdateVideo{
			LessonID: &lesson.Id,
		}

		uid, err := common.FromBase58(*input.VideoID)
		if err != nil {
			return common.NewCustomError(err, "invalid video id", "video_id")
		}

		if err := biz.repo.UpdateVideo(ctx, uid.GetLocalID(), videoUpdate); err != nil {
			return common.ErrCannotUpdateEntity(models.VideoEntityName, err)
		}
	}

	return nil
}
