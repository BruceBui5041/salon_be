package userstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"

	"gorm.io/gorm"
)

func (s *sqlStore) UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	if err := tx.Save(user).Error; err != nil {
		return common.ErrDB(err)
	}

	// Update roles
	if len(user.Roles) > 0 {
		if err := tx.Model(user).Association("Roles").Replace(user.Roles); err != nil {
			return common.ErrDB(err)
		}
	}

	return nil
}

func (s *sqlStore) Update(ctx context.Context, userID uint32, data *models.User) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(data).Error
}
