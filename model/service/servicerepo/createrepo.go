package servicerepo

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"

	"go.uber.org/zap"
)

type ServiceStore interface {
	CreateNewService(ctx context.Context, data *models.Service) error
	Update(ctx context.Context, serviceID uint32, data *models.Service) error
}

type ServiceVersionStore interface {
	CreateNewServiceVersion(ctx context.Context, data *models.ServiceVersion) error
}

type ImageRepo interface {
	CreateImage(ctx context.Context, file *multipart.FileHeader, serviceID uint32, userID uint32) (*models.Image, error)
}

type createServiceRepo struct {
	serviceStore        ServiceStore
	serviceVersionStore ServiceVersionStore
	imageRepo           ImageRepo // Add this field
}

func NewCreateServiceRepo(
	serviceStore ServiceStore,
	serviceVersionStore ServiceVersionStore,
	imageRepo ImageRepo,
) *createServiceRepo {
	return &createServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
		imageRepo:           imageRepo,
	}
}

func (repo *createServiceRepo) CreateNewService(
	ctx context.Context,
	input *servicemodel.CreateService,
) (*models.Service, error) {
	service := &models.Service{
		SQLModel:  common.SQLModel{Status: common.StatusInactive},
		CreatorID: input.CreatorID,
		Slug:      input.Slug,
	}

	if err := repo.serviceStore.CreateNewService(ctx, service); err != nil {
		return nil, common.ErrDB(err)
	}

	if input.ServiceVersion != nil {
		categoryUID, err := common.FromBase58(input.ServiceVersion.CategoryID)
		if err != nil {
			return nil, common.ErrInvalidRequest(errors.New("invalid category ID"))
		}

		subCategoryUID, err := common.FromBase58(input.ServiceVersion.SubCategoryID)
		if err != nil {
			return nil, common.ErrInvalidRequest(errors.New("invalid sub category ID"))
		}

		// Parse IntroVideoID into UID if it's not nil
		var introVideoUID *uint32
		if input.ServiceVersion.IntroVideoID != nil {
			uid, err := common.DecomposeUID(*input.ServiceVersion.IntroVideoID)
			if err != nil {
				return nil, common.ErrInvalidRequest(errors.New("invalid intro video ID"))
			}
			localId := uid.GetLocalID()
			introVideoUID = &localId
		}

		serviceVersion := &models.ServiceVersion{
			SQLModel:      common.SQLModel{Status: common.StatusInactive},
			ServiceID:     service.Id,
			Title:         input.ServiceVersion.Title,
			Description:   input.ServiceVersion.Description,
			CategoryID:    categoryUID.GetLocalID(),
			SubCategoryID: subCategoryUID.GetLocalID(),
			IntroVideoID:  introVideoUID,
			Thumbnail:     input.ServiceVersion.Thumbnail,
			Price:         input.ServiceVersion.Price.GetDecimal(),
			Duration:      input.ServiceVersion.Duration,
		}

		if input.ServiceVersion.DiscountedPrice != nil {
			discounted := input.ServiceVersion.DiscountedPrice.GetDecimal()
			serviceVersion.DiscountedPrice = &discounted
		}

		if err := repo.serviceVersionStore.CreateNewServiceVersion(ctx, serviceVersion); err != nil {
			return nil, common.ErrDB(err)
		}

		service.ServiceVersionID = &serviceVersion.Id
		service.ServiceVersion = serviceVersion

		if len(input.ServiceVersion.Images) > 0 {
			for _, file := range input.ServiceVersion.Images {
				img, err := repo.imageRepo.CreateImage(ctx, file, service.Id, input.CreatorID)
				if err != nil {
					logger.AppLogger.Error(ctx, "faild to upload service image", zap.Error(err))
					return nil, common.ErrDB(err)
				}

				serviceVersion.Images = append(serviceVersion.Images, img)
			}
		}

		if err := repo.serviceStore.Update(ctx, service.Id, service); err != nil {
			logger.AppLogger.Error(ctx, "failed to update service", zap.Error(err))
			return nil, common.ErrDB(err)
		}

		service.ServiceVersionID = &serviceVersion.Id
	}

	return service, nil
}
