package permissionrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/permission/permissionmodel"
)

type CreatePermissionStore interface {
	Create(
		ctx context.Context,
		newPermission *models.Permission,
	) (*models.Permission, error)
}

type createPermissionRepo struct {
	store CreatePermissionStore
}

func NewCreatePermissionRepo(store CreatePermissionStore) *createPermissionRepo {
	return &createPermissionRepo{
		store: store,
	}
}

func (repo *createPermissionRepo) CreateNewPermission(
	ctx context.Context,
	input *permissionmodel.CreatePermission,
) (*models.Permission, error) {
	newPermission := &models.Permission{
		Name:        input.Name,
		Description: input.Description,
		Code:        input.Code,
	}

	permission, err := repo.store.Create(ctx, newPermission)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return permission, nil
}
