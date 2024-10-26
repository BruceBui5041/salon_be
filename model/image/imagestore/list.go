package imagestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) ListImagesByServiceVersionID(
	ctx context.Context,
	serviceVersionID uint32,
	moreKeys ...string,
) ([]models.Image, error) {
	var result []models.Image
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Where("service_version_id = ?", serviceVersionID).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
