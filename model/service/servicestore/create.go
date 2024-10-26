package servicestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) CreateNewService(
	ctx context.Context,
	data *models.Service,
) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
