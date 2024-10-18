package lessonbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
)

type DeleteLessonRepo interface {
	DeleteLesson(ctx context.Context, lessonId uint32) error
	FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error)
	FindLesson(ctx context.Context, conditions map[string]interface{}) (*models.Lesson, error)
}

type deleteLessonBiz struct {
	repo DeleteLessonRepo
}

func NewDeleteLessonBiz(repo DeleteLessonRepo) *deleteLessonBiz {
	return &deleteLessonBiz{repo: repo}
}

func (biz *deleteLessonBiz) DeleteLesson(ctx context.Context, lessonId string) error {
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

	requester := ctx.Value(common.CurrentUser).(common.Requester)
	if course.CreatorID != requester.GetUserId() {
		return common.ErrNoPermission(errors.New("you don't have permission to delete this lesson"))
	}

	if err := biz.repo.DeleteLesson(ctx, lesson.Id); err != nil {
		return common.ErrCannotDeleteEntity(models.LessonEntityName, err)
	}

	return nil
}
