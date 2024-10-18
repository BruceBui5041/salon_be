package commentstore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Comment, error) {
	var comment models.Comment
	db := s.db

	for _, info := range moreInfo {
		db = db.Preload(info)
	}

	if err := db.Where(conditions).First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &comment, nil
}
