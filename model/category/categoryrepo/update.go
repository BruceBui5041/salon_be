package categoryrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/category/categorymodel"

	"github.com/jinzhu/copier"
)

type UpdateCategoryStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
	Update(
		ctx context.Context,
		id uint32,
		data *models.Category,
	) error
}

type updateCategoryRepo struct {
	store UpdateCategoryStore
}

func NewUpdateCategoryRepo(store UpdateCategoryStore) *updateCategoryRepo {
	return &updateCategoryRepo{
		store: store,
	}
}

func (repo *updateCategoryRepo) GetCategory(ctx context.Context, id uint32) (*models.Category, error) {
	category, err := repo.store.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return category, nil
}

func (repo *updateCategoryRepo) UpdateCategory(ctx context.Context, id uint32, data *categorymodel.UpdateCategory) error {
	var categ models.Category
	copier.Copy(&categ, data)

	if err := repo.store.Update(ctx, id, &categ); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
