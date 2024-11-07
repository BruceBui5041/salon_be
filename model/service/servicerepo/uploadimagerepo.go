package servicerepo

import (
	"context"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UploadImagesStore interface {
	Create(ctx context.Context, data *models.Image) error
	List(
		ctx context.Context,
		conditions []interface{},
		moreKeys ...string,
	) ([]*models.Image, error)
}

type M2MServiceVersionImageStore interface {
	Create(ctx context.Context, data *models.M2MServiceVersionImage) error
	Delete(ctx context.Context, conditions map[string]interface{}) error
	List(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]models.M2MServiceVersionImage, error)
}

type ImageUploader interface {
	CreateImage(ctx context.Context, file *multipart.FileHeader, serviceID uint32, userID uint32) (*models.Image, error)
}

type UploadImageServiceVersionStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.ServiceVersion, error)
}

type uploadImagesRepo struct {
	store               UploadImagesStore
	imageUploader       ImageUploader
	serviceVersionStore UploadImageServiceVersionStore
	m2mStore            M2MServiceVersionImageStore
	db                  *gorm.DB
}

func NewUploadImagesRepo(
	store UploadImagesStore,
	imageUploader ImageUploader,
	serviceVersionStore UploadImageServiceVersionStore,
	m2mStore M2MServiceVersionImageStore,
	db *gorm.DB,
) *uploadImagesRepo {
	return &uploadImagesRepo{
		store:               store,
		imageUploader:       imageUploader,
		serviceVersionStore: serviceVersionStore,
		m2mStore:            m2mStore,
		db:                  db,
	}
}

func (repo *uploadImagesRepo) UploadImages(ctx context.Context, data *servicemodel.UploadImages) error {
	serviceID, err := data.GetServiceIDLocalId()
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	var serviceVersion *models.ServiceVersion
	if data.ServiceVersionID != nil {
		serviceVersionID, err := data.GetServiceVersionIDLocalId()
		if err != nil {
			return common.ErrInvalidRequest(err)
		}

		serviceVersion, err = repo.serviceVersionStore.FindOne(
			ctx,
			map[string]interface{}{"id": serviceVersionID},
		)
		if err != nil {
			return common.ErrEntityNotFound(models.ServiceVersionEntityName, err)
		}
	}

	// Use transaction to ensure data consistency
	return repo.db.Transaction(func(tx *gorm.DB) error {
		for _, file := range data.Images {
			// Create image record
			img, err := repo.imageUploader.CreateImage(ctx, file, serviceID, data.UploadedBy)
			if err != nil {
				logger.AppLogger.Error(ctx, "failed to upload service image", zap.Error(err))
				return common.ErrDB(err)
			}

			if err := repo.store.Create(ctx, img); err != nil {
				logger.AppLogger.Error(ctx, "failed to create image record", zap.Error(err))
				return common.ErrDB(err)
			}

			// If service version is specified, create the M2M relationship using the store
			if serviceVersion != nil {
				m2m := &models.M2MServiceVersionImage{
					ServiceVersionID: serviceVersion.Id,
					ImageID:          img.Id,
				}

				if err := repo.m2mStore.Create(ctx, m2m); err != nil {
					logger.AppLogger.Error(ctx, "failed to create m2m relationship", zap.Error(err))
					return common.ErrDB(err)
				}
			}
		}

		return nil
	})
}
