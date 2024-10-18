package userrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"gorm.io/gorm"
)

type UpdateUserStore interface {
	UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error
}

type updateUserRepo struct {
	store UpdateUserStore
}

func NewUpdateUserRepo(store UpdateUserStore) *updateUserRepo {
	return &updateUserRepo{store: store}
}

func (repo *updateUserRepo) UpdateUser(ctx context.Context, tx *gorm.DB, id uint32, data *usermodel.UserUpdate) error {
	var user models.User
	if err := tx.Preload("Roles").First(&user, id).Error; err != nil {
		return common.ErrEntityNotFound(models.UserEntityName, err)
	}

	// Update roles if provided
	if len(data.RoleIDs) > 0 {
		user.Roles = make([]*models.Role, 0, len(data.RoleIDs))
		for _, roleIDStr := range data.RoleIDs {
			roleID, err := common.FromBase58(roleIDStr)
			if err != nil {
				return common.ErrInvalidRequest(err)
			}
			user.Roles = append(user.Roles, &models.Role{SQLModel: common.SQLModel{Id: roleID.GetLocalID()}})
		}
	}

	return repo.store.UpdateUser(ctx, tx, &user)
}
