package commissionstore

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

func (s *sqlStore) Create(ctx context.Context, data *models.Commission) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) Update(ctx context.Context, id uint32, data *models.Commission) error {
	if err := s.db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Commission, error) {
	var commission models.Commission
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).First(&commission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &commission, nil
}
