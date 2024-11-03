package m2mserviceversionimagestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) List(
	ctx context.Context,
	conditions map[string]interface{},
	moreKeys ...string,
) ([]models.M2MServiceVersionImage, error) {
	var result []models.M2MServiceVersionImage
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Where(conditions).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
