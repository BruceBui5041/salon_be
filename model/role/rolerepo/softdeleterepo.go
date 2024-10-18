package rolerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type DeleteRoleStore interface {
	Find(ctx context.Context, cond map[string]interface{}) (*models.Role, error)
	SoftDelete(ctx context.Context, id uint32) error
}

type deleteRoleRepo struct {
	store DeleteRoleStore
}

func NewDeleteRoleRepo(store DeleteRoleStore) *deleteRoleRepo {
	return &deleteRoleRepo{store: store}
}

func (repo *deleteRoleRepo) Find(ctx context.Context, id uint32) (*models.Role, error) {
	role, err := repo.store.Find(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.RoleEntityName, err)
	}

	return role, nil
}

func (repo *deleteRoleRepo) SoftDeleteRole(ctx context.Context, id uint32) error {
	if err := repo.store.SoftDelete(ctx, id); err != nil {
		return common.ErrCannotDeleteEntity(models.RoleEntityName, err)
	}

	return nil
}
