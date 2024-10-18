package categoryrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/category/categorymodel"
)

type CreateCategoryStore interface {
	Create(
		ctx context.Context,
		newCategory *models.Category,
	) (*models.Category, error)
}

type createCategoryRepo struct {
	store CreateCategoryStore
}

func NewCreateCategoryRepo(store CreateCategoryStore) *createCategoryRepo {
	return &createCategoryRepo{
		store: store,
	}
}

func (repo *createCategoryRepo) CreateNewCategory(
	ctx context.Context,
	input *categorymodel.CreateCategory,
) (*models.Category, error) {
	newCategory := &models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	category, err := repo.store.Create(ctx, newCategory)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return category, nil
}
