package rolebiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type DeleteRoleRepo interface {
	Find(ctx context.Context, id uint32) (*models.Role, error)
	SoftDeleteRole(ctx context.Context, id uint32) error
}

type deleteRoleBiz struct {
	repo DeleteRoleRepo
}

func NewDeleteRoleBiz(repo DeleteRoleRepo) *deleteRoleBiz {
	return &deleteRoleBiz{repo: repo}
}

func (biz *deleteRoleBiz) SoftDeleteRole(ctx context.Context, id uint32) error {
	oldData, err := biz.repo.Find(ctx, id)
	if err != nil {
		return common.ErrCannotGetEntity(models.RoleEntityName, err)
	}

	if oldData.Status == common.StatusInactive {
		return common.ErrEntityDeleted(models.RoleEntityName, nil)
	}

	if err := biz.repo.SoftDeleteRole(ctx, id); err != nil {
		return common.ErrCannotDeleteEntity(models.RoleEntityName, err)
	}

	return nil
}
