package notificationstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	data *models.Notification,
) error {
	if err := s.db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
