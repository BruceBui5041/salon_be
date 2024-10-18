package commentstore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	updateData *models.Comment,
) error {
	if err := s.db.Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
