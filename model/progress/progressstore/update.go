package progressstore

import (
	"context"
	"video_server/common"
	models "video_server/model"
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
