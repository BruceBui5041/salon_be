package categorybiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/category/categorymodel"
)

type CategoryRepo interface {
	CreateNewCategory(ctx context.Context, input *categorymodel.CreateCategory) (*models.Category, error)
}

type createCategoryBiz struct {
	repo CategoryRepo
}

func NewCreateCategoryBiz(repo CategoryRepo) *createCategoryBiz {
	return &createCategoryBiz{repo: repo}
}

func (c *createCategoryBiz) CreateNewCategory(ctx context.Context, input *categorymodel.CreateCategory) error {
	if input.Name == "" {
		return common.ErrInvalidRequest(errors.New("category name is required"))
	}

	if len(input.Name) > 100 {
		return common.ErrInvalidRequest(errors.New("category name must not exceed 100 characters"))
	}

	if input.Code == "" {
		return common.ErrInvalidRequest(errors.New("code is required"))
	}

	if len(input.Code) > 100 {
		return common.ErrInvalidRequest(errors.New("code must not exceed 100 characters"))
	}

	category, err := c.repo.CreateNewCategory(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.CategoryEntityName, err)
	}

	input.Id = category.Id

	return nil
}
