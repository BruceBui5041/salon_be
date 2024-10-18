package progressstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newProgress *models.Progress,
) (*models.Progress, error) {
	if err := s.db.Create(newProgress).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return newProgress, nil
}
