package lecturestore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newLecture *models.Lecture,
) (uint32, error) {
	if err := s.db.Create(newLecture).Error; err != nil {
		return 0, err
	}

	return newLecture.Id, nil
}
