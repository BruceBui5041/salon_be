package coursebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
)

type CourseRepo interface {
	CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) (*models.Course, error)
}

type createCourseBiz struct {
	repo CourseRepo
}

func NewCreateCourseBiz(repo CourseRepo) *createCourseBiz {
	return &createCourseBiz{repo: repo}
}

func (c *createCourseBiz) CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) (*models.Course, error) {
	// Validate required fields
	if input.Title == "" {
		return nil, errors.New("course title is required")
	}

	if input.CreatorID == 0 {
		return nil, errors.New("creator ID is required")
	}

	if input.CategoryID == "" {
		return nil, errors.New("category ID is required")
	}

	if input.Slug == "" {
		return nil, errors.New("course slug is required")
	}

	if len(input.Title) > 255 {
		return nil, errors.New("course title must not exceed 255 characters")
	}

	if len(input.Slug) > 255 {
		return nil, errors.New("course slug must not exceed 255 characters")
	}

	// Call the storage layer to create the course
	course, err := c.repo.CreateNewCourse(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.CourseEntityName, err)
	}

	return course, nil
}
