// userdevice/userdevicestore/create.go
package userdevicestore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newUserDevice *models.UserDevice,
) (*models.UserDevice, error) {
	if err := s.db.Create(newUserDevice).Error; err != nil {
		return nil, err
	}
	return newUserDevice, nil
}

func (s *sqlStore) Delete(
	ctx context.Context,
	id uint32,
) error {
	if err := s.db.Where("id = ?", id).Delete(&models.UserDevice{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) Find(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.UserDevice, error) {
	var userDevices []models.UserDevice
	db := s.db
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	if err := db.Where(conditions).Find(&userDevices).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return userDevices, nil
}

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.UserDevice, error) {
	var userDevice *models.UserDevice
	db := s.db
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	if err := db.Where(conditions).First(&userDevice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return userDevice, nil
}

type sqlStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	data *models.UserDevice,
) error {
	if err := s.db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
