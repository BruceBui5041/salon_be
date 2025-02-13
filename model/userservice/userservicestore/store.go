package userservicestore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"

	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}

func (s *sqlStore) DeleteByCondition(ctx context.Context, conditions map[string]interface{}) error {
	if err := s.db.Where(conditions).Delete(&models.UserService{}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) Create(ctx context.Context, data *models.UserService) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
