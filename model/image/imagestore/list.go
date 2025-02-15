package imagestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) List(
	ctx context.Context,
	conditions []interface{},
	moreKeys ...string,
) ([]*models.Image, error) {
	var result []*models.Image
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Find(&result, conditions...).Error; err != nil {
		return nil, err
	}

	return result, nil
}
