package permissionrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/permission/permissionmodel"

	"github.com/jinzhu/copier"
)

type UpdatePermissionStore interface {
	Update(ctx context.Context, id uint32, data *models.Permission) error
}

type updatePermissionRepo struct {
	store UpdatePermissionStore
}

func NewUpdatePermissionRepo(store UpdatePermissionStore) *updatePermissionRepo {
	return &updatePermissionRepo{store: store}
}

func (repo *updatePermissionRepo) UpdatePermission(ctx context.Context, id uint32, input *permissionmodel.UpdatePermission) error {
	permission := models.Permission{}

	if err := copier.Copy(&permission, input); err != nil {
		return common.ErrInvalidRequest(err)
	}

	if err := repo.store.Update(ctx, id, &permission); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
