package locationstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreKeys ...string,
) (*models.Location, error) {
	var result models.Location
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
