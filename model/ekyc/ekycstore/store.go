package ekycstore

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

func (s *sqlStore) CreateKYCProfile(ctx context.Context, data *models.KYCProfile) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) CreateDocument(ctx context.Context, data *models.IDDocument) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) CreateFaceVerification(ctx context.Context, data *models.FaceVerification) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) CreateKYCImage(ctx context.Context, data *models.KYCImageUpload) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) UpdateKYCImage(ctx context.Context, kycImageId uint32, data *models.KYCImageUpload) error {
	if err := s.db.Where("id = ?", kycImageId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
