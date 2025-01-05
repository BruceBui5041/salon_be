package bookingstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Create(ctx context.Context, data *models.Booking) (uint32, error) {
	if err := s.db.Create(data).Error; err != nil {
		return 0, err
	}
	return data.Id, nil
}
