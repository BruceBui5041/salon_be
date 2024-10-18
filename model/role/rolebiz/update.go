package rolebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/role/rolemodel"
)

type UpdateRoleRepo interface {
	GetRole(ctx context.Context, id uint32) (*models.Role, error)
	UpdateRole(ctx context.Context, id uint32, input *rolemodel.UpdateRole) error
}

type updateRoleBiz struct {
	repo UpdateRoleRepo
}

func NewUpdateRoleBiz(repo UpdateRoleRepo) *updateRoleBiz {
	return &updateRoleBiz{repo: repo}
}

func (biz *updateRoleBiz) UpdateRole(ctx context.Context, id uint32, input *rolemodel.UpdateRole) error {
	if input.Name == "" {
		return common.ErrInvalidRequest(errors.New("role name is required"))
	}

	if len(input.Name) > 50 {
		return common.ErrInvalidRequest(errors.New("role name must not exceed 50 characters"))
	}

	if input.Code == "" {
		return common.ErrInvalidRequest(errors.New("role code is required"))
	}

	if len(input.Code) > 50 {
		return common.ErrInvalidRequest(errors.New("role code must not exceed 50 characters"))
	}

	oldData, err := biz.repo.GetRole(ctx, id)
	if err != nil {
		return common.ErrCannotGetEntity(models.RoleEntityName, err)
	}

	if oldData.Status == "inactive" {
		return common.ErrEntityDeleted(models.RoleEntityName, nil)
	}

	if err := biz.repo.UpdateRole(ctx, id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.RoleEntityName, err)
	}

	return nil
}
