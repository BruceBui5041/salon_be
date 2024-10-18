package lessonstore

import (
	"context"
	"video_server/model/lesson/lessonmodel"
)

func (s *sqlStore) UpdateLesson(
	ctx context.Context,
	lessonId uint32,
	data *lessonmodel.UpdateLesson,
) error {
	if err := s.db.Where("id = ?", lessonId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
