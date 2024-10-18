package permissionbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/permission/permissionmodel"
)

type PermissionRepo interface {
	CreateNewPermission(ctx context.Context, input *permissionmodel.CreatePermission) (*models.Permission, error)
}

type createPermissionBiz struct {
	repo PermissionRepo
}

func NewCreatePermissionBiz(repo PermissionRepo) *createPermissionBiz {
	return &createPermissionBiz{repo: repo}
}

func (c *createPermissionBiz) CreateNewPermission(ctx context.Context, input *permissionmodel.CreatePermission) error {
	if input.Name == "" {
		return common.ErrInvalidRequest(errors.New("permission name is required"))
	}

	if len(input.Name) > 50 {
		return common.ErrInvalidRequest(errors.New("permission name must not exceed 50 characters"))
	}

	if input.Code == "" {
		return common.ErrInvalidRequest(errors.New("permission code is required"))
	}

	if len(input.Code) > 50 {
		return common.ErrInvalidRequest(errors.New("permission code must not exceed 50 characters"))
	}

	_, err := c.repo.CreateNewPermission(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.PermissionEntityName, err)
	}

	return nil
}
