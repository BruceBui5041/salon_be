package userprofilestore

import (
	"context"
	models "salon_be/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.UserProfile, error) {
	var userProfile models.UserProfile

	db := s.db.Table(models.UserProfile{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &userProfile, nil
}
