package userstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.User, error) {
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	user := &models.User{}

	if err := db.Where(conditions).First(&user).Error; err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}

	return user, nil
}
