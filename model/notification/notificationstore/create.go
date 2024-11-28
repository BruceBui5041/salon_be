package notificationstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newNotification *models.Notification,
) (*models.Notification, error) {
	if err := s.db.Create(newNotification).Error; err != nil {
		return nil, err
	}

	return newNotification, nil
}
