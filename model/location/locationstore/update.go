package locationstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(ctx context.Context, data *models.Location) error {
	return s.db.Model(&models.Location{}).
		Where("id = ?", data.Id).
		Updates(data).Error
}
