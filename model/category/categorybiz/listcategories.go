package categorybiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/category/categorymodel"

	"github.com/jinzhu/copier"
)

type CategoryStore interface {
	FindAll(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Category, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type categoryBiz struct {
	categoryStore CategoryStore
}

func NewCategoryBiz(categoryStore CategoryStore) *categoryBiz {
	return &categoryBiz{categoryStore: categoryStore}
}

func (biz *categoryBiz) ListCategories(ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]*categorymodel.CategoryResponse, error) {
	extentConds := conditions
	extentConds["status"] = "active"
	categories, err := biz.categoryStore.FindAll(ctx, conditions, moreInfo...)
	if err != nil {
		return nil, common.ErrCannotListEntity(models.CategoryEntityName, err)
	}

	var categoriesRes []*categorymodel.CategoryResponse
	err = copier.Copy(&categoriesRes, categories)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	for _, categ := range categoriesRes {
		// categ.Mask(false)
		categ.CountService()
		categ.RemoveServicesResponse()
	}

	return categoriesRes, nil
}
