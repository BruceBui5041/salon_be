// File: lecturebiz/updatelecture.go

package lecturebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lecture/lecturemodel"
)

type UpdateLectureRepo interface {
	UpdateLecture(ctx context.Context, lectureId uint32, input *lecturemodel.UpdateLecture) error
	FindCourse(ctx context.Context, conditions map[string]interface{}) (*models.Course, error)
	FindLecture(ctx context.Context, conditions map[string]interface{}) (*models.Lecture, error)
}

type updateLectureBiz struct {
	repo UpdateLectureRepo
}

func NewUpdateLectureBiz(repo UpdateLectureRepo) *updateLectureBiz {
	return &updateLectureBiz{repo: repo}
}

func (biz *updateLectureBiz) UpdateLecture(ctx context.Context, lectureId string, input *lecturemodel.UpdateLecture) error {
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
		return common.ErrNoPermission(errors.New("you don't have permission to update this lecture"))
	}

	if err := biz.repo.UpdateLecture(ctx, lecture.Id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.LectureEntityName, err)
	}

	return nil
}
