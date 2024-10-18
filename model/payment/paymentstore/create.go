package paymentstore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newPayment *models.Payment,
) (uint32, error) {
	if err := s.db.Create(newPayment).Error; err != nil {
		return 0, err
	}

	return newPayment.Id, nil
}
