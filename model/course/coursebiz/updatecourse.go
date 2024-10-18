package coursebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
)

type UpdateCourseRepo interface {
	UpdateCourse(ctx context.Context, id uint32, input *coursemodel.UpdateCourse) error
}

type updateCourseBiz struct {
	repo UpdateCourseRepo
}

func NewUpdateCourseBiz(repo UpdateCourseRepo) *updateCourseBiz {
	return &updateCourseBiz{repo: repo}
}

func (c *updateCourseBiz) UpdateCourse(ctx context.Context, id uint32, input *coursemodel.UpdateCourse) error {
	if input.Title != "" && len(input.Title) > 255 {
		return errors.New("course title must not exceed 255 characters")
	}

	err := c.repo.UpdateCourse(ctx, id, input)
	if err != nil {
		return common.ErrCannotUpdateEntity(models.CourseEntityName, err)
	}

	return nil
}
