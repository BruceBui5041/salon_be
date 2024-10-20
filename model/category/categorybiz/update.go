package categorybiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/category/categorymodel"
)

type UpdateCategoryRepo interface {
	GetCategory(ctx context.Context, id uint32) (*models.Category, error)
	UpdateCategory(ctx context.Context, id uint32, data *categorymodel.UpdateCategory) error
}

type updateCategoryBiz struct {
	repo UpdateCategoryRepo
}

func NewUpdateCategoryBiz(repo UpdateCategoryRepo) *updateCategoryBiz {
	return &updateCategoryBiz{repo: repo}
}

func (biz *updateCategoryBiz) UpdateCategory(ctx context.Context, id uint32, data *categorymodel.UpdateCategory) error {
	oldData, err := biz.repo.GetCategory(ctx, id)
	if err != nil {
		return common.ErrCannotGetEntity(models.CategoryEntityName, err)
	}

	if oldData == nil {
		return common.ErrEntityNotFound(models.CategoryEntityName, errors.New("category not found"))
	}

	if data.Name != "" && len(data.Name) > 100 {
		return common.ErrInvalidRequest(errors.New("category name must not exceed 100 characters"))
	}

	if data.Code != "" && len(data.Code) > 100 {
		return common.ErrInvalidRequest(errors.New("code must not exceed 100 characters"))
	}

	if err := biz.repo.UpdateCategory(ctx, id, data); err != nil {
		return common.ErrCannotUpdateEntity(models.CategoryEntityName, err)
	}

	return nil
}
