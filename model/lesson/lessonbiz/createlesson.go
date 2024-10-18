package lessonbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"
)

type LessonRepo interface {
	CreateNewLesson(ctx context.Context, input *lessonmodel.CreateLesson) (*models.Lesson, error)
}

type createLessonBiz struct {
	repo LessonRepo
}

func NewCreateLessonBiz(repo LessonRepo) *createLessonBiz {
	return &createLessonBiz{repo: repo}
}

func (c *createLessonBiz) CreateNewLesson(ctx context.Context, input *lessonmodel.CreateLesson) (*models.Lesson, error) {
	if input.Title == "" {
		return nil, errors.New("lesson title is required")
	}

	if input.CourseID == "" {
		return nil, errors.New("course ID is required")
	}

	if len(input.Title) > 255 {
		return nil, errors.New("lesson title must not exceed 255 characters")
	}

	lesson, err := c.repo.CreateNewLesson(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.LessonEntityName, err)
	}

	return lesson, nil
}
