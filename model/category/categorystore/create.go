package categorystore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newCategory *models.Category,
) (*models.Category, error) {
	if err := s.db.Create(newCategory).Error; err != nil {
		return nil, err
	}

	return newCategory, nil
}
