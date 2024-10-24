package videostore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) Exist(
	ctx context.Context,
	conditions map[string]interface{},
) (bool, error) {
	var count int64

	db := s.db.Table(models.Video{}.TableName())

	if err := db.Where(conditions).Count(&count).Error; err != nil {
		return false, common.ErrDB(err)
	}

	return count > 0, nil
}
