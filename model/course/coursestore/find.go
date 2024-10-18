package coursestore

import (
	"context"
	"video_server/common"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindAll(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...interface{},
) ([]models.Course, error) {
	var courses []models.Course
	db := s.db

	for _, info := range moreInfo {
		switch v := info.(type) {
		case string:
			db = db.Preload(v)
		case common.PreloadInfo:
			if v.Function != nil {
				db = db.Preload(v.Name, v.Function)
			} else {
				db = db.Preload(v.Name)
			}
		}
	}

	if err := db.Where(conditions).Find(&courses).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return courses, nil
}
