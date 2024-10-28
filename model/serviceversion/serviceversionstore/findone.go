package serviceversionstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.ServiceVersion, error) {
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	serviceVersion := &models.ServiceVersion{}

	if err := db.Where(conditions).First(&serviceVersion).Error; err != nil {
		return nil, common.ErrEntityNotFound(models.ServiceVersionEntityName, err)
	}

	return serviceVersion, nil
}
