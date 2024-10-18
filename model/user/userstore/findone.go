package userstore

import (
	"context"
	"video_server/common"
	models "video_server/model"
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
		return nil, common.ErrCannotListEntity(models.UserEntityName, err)
	}

	return user, nil
}
