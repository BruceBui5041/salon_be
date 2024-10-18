package rolerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/role/rolemodel"
)

type UpdateRoleStore interface {
	Find(ctx context.Context, cond map[string]interface{}) (*models.Role, error)
	Update(ctx context.Context, id uint32, data *rolemodel.UpdateRole) error
	DeleteRolePermissions(ctx context.Context, roleID uint32) error
	CreateRolePermission(ctx context.Context, rolePermissions []models.RolePermission) error
}

type updateRoleRepo struct {
	store UpdateRoleStore
}

func NewUpdateRoleRepo(store UpdateRoleStore) *updateRoleRepo {
	return &updateRoleRepo{store: store}
}

func (repo *updateRoleRepo) GetRole(ctx context.Context, id uint32) (*models.Role, error) {
	role, err := repo.store.Find(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.RoleEntityName, err)
	}

	return role, nil
}

func (repo *updateRoleRepo) UpdateRole(ctx context.Context, id uint32, input *rolemodel.UpdateRole) error {
	if err := repo.store.Update(ctx, id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.RoleEntityName, err)
	}

	// Delete existing role permissions
	if err := repo.store.DeleteRolePermissions(ctx, id); err != nil {
		return common.ErrCannotDeleteEntity(models.RolePermissionEntityName, err)
	}

	if len(input.PermissionInfo) != 0 {
		rolePermissions := make([]models.RolePermission, 0, len(input.PermissionInfo))
		for _, perm := range input.PermissionInfo {
			permID, err := common.FromBase58(perm.ID)
			if err != nil {
				return common.ErrInvalidRequest(err)
			}

			rolePermissions = append(rolePermissions, models.RolePermission{
				RoleID:           id,
				PermissionID:     permID.GetLocalID(),
				CreatePermission: perm.CreatePermission,
				ReadPermission:   perm.ReadPermission,
				WritePermission:  perm.WritePermission,
				DeletePermission: perm.DeletePermission,
			})
		}

		if err := repo.store.CreateRolePermission(ctx, rolePermissions); err != nil {
			return common.ErrCannotCreateEntity(models.RolePermissionEntityName, err)
		}
	}

	return nil
}
