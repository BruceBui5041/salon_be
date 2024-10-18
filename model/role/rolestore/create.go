package rolestore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Create(ctx context.Context, newRole *models.Role) error {
	if err := s.db.Create(newRole).Error; err != nil {
		return err
	}

	return nil
}

func (s *sqlStore) CreateRolePermission(ctx context.Context, rolePermissions []models.RolePermission) error {
	if err := s.db.Create(&rolePermissions).Error; err != nil {
		return err
	}

	return nil
}
