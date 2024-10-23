package serviceversionstore

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

func (s *sqlStore) CreateNewServiceVersion(
	ctx context.Context,
	data *models.ServiceVersion,
) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
