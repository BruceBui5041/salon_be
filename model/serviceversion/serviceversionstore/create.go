package serviceversionstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) CreateNewServiceVersion(
	ctx context.Context,
	data *models.ServiceVersion,
) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
