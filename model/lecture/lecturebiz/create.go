package lecturebiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lecture/lecturemodel"
)

type LectureRepo interface {
	CreateNewLecture(ctx context.Context, input *lecturemodel.CreateLecture) (*models.Lecture, error)
}

type createLectureBiz struct {
	repo LectureRepo
}

func NewCreateLectureBiz(repo LectureRepo) *createLectureBiz {
	return &createLectureBiz{repo: repo}
}

func (c *createLectureBiz) CreateNewLecture(ctx context.Context, input *lecturemodel.CreateLecture) (*models.Lecture, error) {
	if input.Title == "" {
		return nil, errors.New("lecture title is required")
	}

	if input.CourseID == "" {
		return nil, errors.New("course ID is required")
	}

	if len(input.Title) > 255 {
		return nil, errors.New("lecture title must not exceed 255 characters")
	}

	lecture, err := c.repo.CreateNewLecture(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.LectureEntityName, err)
	}

	return lecture, nil
}
