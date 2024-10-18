package videostore

import (
	"context"
	"salon_be/model/video/videomodel"
)

func (s *sqlStore) UpdateVideo(
	ctx context.Context,
	id uint32,
	updateData *videomodel.UpdateVideo,
) error {
	db := s.db

	// Update the current video with the provided updateData
	if err := db.Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
