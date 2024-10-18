package lecturestore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Lecture, error) {
	var lecture models.Lecture
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := s.db.Where(conditions).First(&lecture).Error; err != nil {
		return nil, err
	}

	return &lecture, nil
}
