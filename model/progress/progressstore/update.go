package progressstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	conditions map[string]interface{},
	updateData *models.Progress,
) error {
	if err := s.db.Where(conditions).Updates(updateData).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
