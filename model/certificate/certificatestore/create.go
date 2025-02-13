package certificatestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Create(ctx context.Context, data *models.Certificate) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
