package videoprocessinfostore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Find(
	ctx context.Context,
	conditions map[string]interface{},
	query string,
	moreInfo ...string,
) ([]*models.VideoProcessInfo, error) {
	var processInfo []*models.VideoProcessInfo
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if query != "" {
		db = db.Where(query)
	} else {
		db = db.Where(conditions)
	}

	if err := db.Find(&processInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return processInfo, nil
}
