package groupproviderstore

import (
	"context"
	models "salon_be/model"

	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}

func (s *sqlStore) Create(ctx context.Context, data *models.GroupProvider) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreKeys ...string,
) (*models.GroupProvider, error) {
	var result models.GroupProvider
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *sqlStore) Update(ctx context.Context, id uint32, data *models.GroupProvider) error {
	if err := s.db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	if len(data.Admins) > 0 {
		if err := s.db.Model(data).Association("Admins").Replace(data.Admins); err != nil {
			return err
		}
	}

	return nil
}
