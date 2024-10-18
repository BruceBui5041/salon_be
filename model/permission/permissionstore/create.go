package permissionstore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newPermission *models.Permission,
) (*models.Permission, error) {
	if err := s.db.Create(newPermission).Error; err != nil {
		return nil, err
	}

	return newPermission, nil
}
