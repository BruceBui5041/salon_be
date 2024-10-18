package rolerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/role/rolemodel"
)

type CreateRoleStore interface {
	Create(ctx context.Context, newRole *models.Role) error
	CreateRolePermission(ctx context.Context, rolePermissions []models.RolePermission) error
}

type createRoleRepo struct {
	store CreateRoleStore
}

func NewCreateRoleRepo(store CreateRoleStore) *createRoleRepo {
	return &createRoleRepo{store: store}
}

func (repo *createRoleRepo) CreateNewRole(ctx context.Context, input *rolemodel.CreateRole) (*models.Role, error) {
	newRole := &models.Role{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
	}

	if err := repo.store.Create(ctx, newRole); err != nil {
		return nil, common.ErrDB(err)
	}

	if len(input.Permissions) != 0 {
		rolePermissions := make([]models.RolePermission, 0, len(input.Permissions))
		for _, perm := range input.Permissions {
			permID, err := common.FromBase58(perm.ID)
			if err != nil {
				return nil, common.ErrInvalidRequest(err)
			}

			rolePermissions = append(rolePermissions, models.RolePermission{
				RoleID:           newRole.Id,
				PermissionID:     permID.GetLocalID(),
				CreatePermission: perm.CreatePermission,
				ReadPermission:   perm.ReadPermission,
				WritePermission:  perm.WritePermission,
				DeletePermission: perm.DeletePermission,
			})
		}

		if err := repo.store.CreateRolePermission(ctx, rolePermissions); err != nil {
			return nil, common.ErrDB(err)
		}
	}

	return newRole, nil
}
