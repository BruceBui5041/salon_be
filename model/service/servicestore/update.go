package servicestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(ctx context.Context, serviceID uint32, updates *models.Service) error {
	return s.db.Model(&models.Service{}).
		Where("id = ?", serviceID).
		Updates(updates).Error
}
