package lessonstore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) CreateNewLesson(
	ctx context.Context,
	newLesson *models.Lesson,
) (uint32, error) {
	if err := s.db.Create(newLesson).Error; err != nil {
		return 0, err
	}

	return newLesson.Id, nil
}
