package servicestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) UpdateMany(ctx context.Context, condition map[string]interface{}, data *models.Service) error {
	return s.db.Model(&models.Service{}).
		Where(condition).
		Updates(data).Error
}
