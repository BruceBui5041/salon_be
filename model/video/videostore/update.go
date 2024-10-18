package videostore

import (
	"context"
	"video_server/model/video/videomodel"
)

func (s *sqlStore) UpdateVideo(
	ctx context.Context,
	id uint32,
	updateData *videomodel.UpdateVideo,
) error {
	db := s.db

	// Check if the LessonID is provided in the updateData
	if updateData.LessonID != nil {
		// Update all existing videos with the same LessonID to NULL
		if err := db.Model(&videomodel.UpdateVideo{}).
			Where("lesson_id = ?", *updateData.LessonID).
			Update("lesson_id", nil).Error; err != nil {
			return err
		}
	}

	// Update the current video with the provided updateData
	if err := db.Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
