package servicerepo

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
	"salon_be/utils/customtypes"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UpdateServiceStore interface {
	Update(ctx context.Context, serviceID uint32, data *models.Service) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Service, error)
}

type UpdateServiceVersionStore interface {
	Update(ctx context.Context, versionID uint32, data *models.ServiceVersion) error
}

type UpdateImageRepo interface {
	CreateImage(ctx context.Context, file *multipart.FileHeader, serviceID uint32, userID uint32) (*models.Image, error)
}

type updateServiceRepo struct {
	serviceStore        UpdateServiceStore
	serviceVersionStore UpdateServiceVersionStore
	imageRepo           UpdateImageRepo
}

func NewUpdateServiceRepo(
	serviceStore UpdateServiceStore,
	serviceVersionStore UpdateServiceVersionStore,
	imageRepo UpdateImageRepo,
) *updateServiceRepo {
	return &updateServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
		imageRepo:           imageRepo,
	}
}

func (repo *updateServiceRepo) UpdateService(
	ctx context.Context,
	input *servicemodel.UpdateService,
) (*models.Service, error) {
	serviceVersionId, err := input.GetServiceVersionLocalId()
	if err != nil {
		return nil, err
	}

	serviceID, err := input.GetServiceLocalId()
	if err != nil {
		return nil, err
	}

	// Find existing service to get creator_id
	existingService, err := repo.serviceStore.FindOne(ctx, map[string]interface{}{"id": serviceID})
	if err != nil {
		return nil, common.ErrDB(err)
	}

	if existingService == nil {
		return nil, common.ErrEntityNotFound("service", nil)
	}

	service := &models.Service{
		ServiceVersionID: &serviceVersionId,
	}

	if err := repo.serviceStore.Update(ctx, serviceID, service); err != nil {
		return nil, common.ErrDB(err)
	}

	if input.ServiceVersion != nil {
		categoryID, err := input.ServiceVersion.GetCateogryLocalId(ctx)
		if err != nil {
			return nil, err
		}

		subCategoryID, err := input.ServiceVersion.GetSubCategoryLocalId(ctx)
		if err != nil {
			return nil, err
		}

		var introVideoUID *uint32
		if input.ServiceVersion.IntroVideoID != "" {
			introVideoID, err := input.ServiceVersion.GetIntroVideoLocalId(ctx)
			if err != nil {
				return nil, err
			}
			introVideoUID = &introVideoID
		}

		serviceVersion := &models.ServiceVersion{
			ServiceID:     serviceID,
			Title:         input.ServiceVersion.Title,
			Description:   input.ServiceVersion.Description,
			CategoryID:    categoryID,
			SubCategoryID: subCategoryID,
			IntroVideoID:  introVideoUID,
			Thumbnail:     input.ServiceVersion.Thumbnail,
			Price:         input.ServiceVersion.Price.GetDecimal(),
			Duration:      input.ServiceVersion.Duration,
		}

		if input.ServiceVersion.Price != customtypes.DecimalString(decimal.Zero) {
			price := input.ServiceVersion.Price.GetDecimal()
			if price.Equal(decimal.Zero) {
				return nil, common.ErrInvalidRequest(errors.New("price must be greater than 0"))
			}
		}

		if input.ServiceVersion.DiscountedPrice != nil {
			discounted := input.ServiceVersion.DiscountedPrice.GetDecimal()
			if discounted.Decimal.GreaterThanOrEqual(serviceVersion.Price) {
				return nil, common.ErrInvalidRequest(errors.New("discount price must be less than price"))
			}
			serviceVersion.DiscountedPrice = &discounted
		}

		if input.ServiceVersion.Duration < 900 {
			return nil, common.ErrInvalidRequest(errors.New("duration must be at least 15 minutes"))
		}

		// Handle image uploads
		if len(input.ServiceVersion.Images) > 0 {
			for _, file := range input.ServiceVersion.Images {
				img, err := repo.imageRepo.CreateImage(ctx, file, serviceID, existingService.CreatorID)
				if err != nil {
					logger.AppLogger.Error(ctx, "failed to upload service image", zap.Error(err))
					return nil, common.ErrDB(err)
				}

				serviceVersion.Images = append(serviceVersion.Images, img)
			}
		}

		if err := repo.serviceVersionStore.Update(ctx, serviceVersionId, serviceVersion); err != nil {
			return nil, common.ErrDB(err)
		}

		service.ServiceVersionID = &serviceVersion.Id
		service.ServiceVersion = serviceVersion
	}

	return service, nil
}

func (repo *updateServiceRepo) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Service, error) {
	return repo.serviceStore.FindOne(ctx, conditions, moreInfo...)
}
