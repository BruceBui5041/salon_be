package lecturebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
)

type DeleteLectureRepo interface {
	DeleteLecture(ctx context.Context, lectureId uint32) error
	FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error)
	FindLecture(ctx context.Context, conditions map[string]interface{}) (*models.Lecture, error)
}

type deleteLectureBiz struct {
	repo DeleteLectureRepo
}

func NewDeleteLectureBiz(repo DeleteLectureRepo) *deleteLectureBiz {
	return &deleteLectureBiz{repo: repo}
}

func (biz *deleteLectureBiz) DeleteLecture(ctx context.Context, lectureId string) error {
	lectureUid, err := common.FromBase58(lectureId)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	lecture, err := biz.repo.FindLecture(ctx, map[string]interface{}{"id": lectureUid.GetLocalID()})
	if err != nil {
		return common.ErrCannotGetEntity(models.LectureEntityName, err)
	}

	if lecture == nil {
		return common.ErrEntityNotFound(models.LectureEntityName, errors.New("lecture not found"))
	}

	course, err := biz.repo.FindCourse(ctx, map[string]interface{}{"id": lecture.CourseID})
	if err != nil {
		return common.ErrCannotGetEntity(models.CourseEntityName, err)
	}

	requester := ctx.Value(common.CurrentUser).(common.Requester)
	if course.CreatorID != requester.GetUserId() {
		return common.ErrNoPermission(errors.New("you don't have permission to delete this lecture"))
	}

	if err := biz.repo.DeleteLecture(ctx, lecture.Id); err != nil {
		return common.ErrCannotDeleteEntity(models.LectureEntityName, err)
	}

	return nil
}
