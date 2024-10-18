package permissionbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/permission/permissionmodel"
)

type UpdatePermissionRepo interface {
	UpdatePermission(ctx context.Context, id uint32, input *permissionmodel.UpdatePermission) error
}

type updatePermissionBiz struct {
	repo UpdatePermissionRepo
}

func NewUpdatePermissionBiz(repo UpdatePermissionRepo) *updatePermissionBiz {
	return &updatePermissionBiz{repo: repo}
}

func (biz *updatePermissionBiz) UpdatePermission(ctx context.Context, id string, input *permissionmodel.UpdatePermission) error {
	if input.Name != nil && *input.Name == "" {
		return common.ErrInvalidRequest(errors.New("permission name cannot be empty"))
	}

	if input.Name != nil && len(*input.Name) > 50 {
		return common.ErrInvalidRequest(errors.New("permission name must not exceed 50 characters"))
	}

	uid, err := common.FromBase58(id)
	if err != nil {
		panic(common.ErrInvalidRequest(err))
	}

	if err := biz.repo.UpdatePermission(ctx, uid.GetLocalID(), input); err != nil {
		return common.ErrCannotUpdateEntity(models.PermissionEntityName, err)
	}

	return nil
}
