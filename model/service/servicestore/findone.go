package servicestore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Service, error) {
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	service := &models.Service{}

	if err := db.Where(conditions).First(&service).Error; err != nil {
		return nil, common.ErrCannotListEntity(models.ServiceEntityName, err)
	}

	return service, nil
}
