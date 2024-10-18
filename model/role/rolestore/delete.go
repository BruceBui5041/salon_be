package rolestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) DeleteRolePermissions(ctx context.Context, roleID uint32) error {
	if err := s.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	return nil
}
