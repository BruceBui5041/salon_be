package notificationdetailstore

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

func (s *sqlStore) Create(
	ctx context.Context,
	detail *models.NotificationDetail,
) (*models.NotificationDetail, error) {
	if err := s.db.Create(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *sqlStore) Delete(
	ctx context.Context,
	conditions map[string]interface{},
) error {
	if err := s.db.Where(conditions).Delete(&models.NotificationDetail{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) Find(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.NotificationDetail, error) {
	var details []models.NotificationDetail
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).Find(&details).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return details, nil
}

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.NotificationDetail, error) {
	var detail *models.NotificationDetail
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).First(&detail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return detail, nil
}

func (s *sqlStore) Update(
	ctx context.Context,
	conditions map[string]interface{},
	data *models.NotificationDetail,
) error {
	if err := s.db.Where(conditions).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
