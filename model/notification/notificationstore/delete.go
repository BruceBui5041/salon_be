package notificationstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Delete(
	ctx context.Context,
	id uint32,
) error {
	if err := s.db.Where("id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		return err
	}

	return nil
}
