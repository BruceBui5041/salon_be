package lessonstore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Lesson, error) {
	var lesson models.Lesson
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := s.db.Where(conditions).First(&lesson).Error; err != nil {
		return nil, err
	}

	return &lesson, nil
}
