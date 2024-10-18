package rolestore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Find(ctx context.Context, cond map[string]interface{}) (*models.Role, error) {
	var role models.Role
	if err := s.db.Where(cond).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
