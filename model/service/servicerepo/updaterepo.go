package servicerepo

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/service/serviceerror"
	"salon_be/model/service/servicemodel"
	"salon_be/utils/customtypes"

	"github.com/samber/lo"
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
	CreateNewServiceVersion(ctx context.Context, data *models.ServiceVersion) error
	Update(ctx context.Context, versionID uint32, data *models.ServiceVersion) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.ServiceVersion, error)
}

type UpdateImageRepo interface {
	CreateImage(ctx context.Context, file *multipart.FileHeader, serviceID uint32, userID uint32) (*models.Image, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.Image, error)
	List(
		ctx context.Context,
		conditions []interface{},
		moreKeys ...string,
	) ([]*models.Image, error)
}

type UpdateM2MVersionImageStore interface {
	Delete(ctx context.Context, conditions map[string]interface{}) error
	Create(ctx context.Context, data *models.M2MServiceVersionImage) error
	Update(ctx context.Context, updates *models.M2MServiceVersionImage) error
}

type updateServiceRepo struct {
	serviceStore        UpdateServiceStore
	serviceVersionStore UpdateServiceVersionStore
	imageRepo           UpdateImageRepo
	m2mVersionImages    UpdateM2MVersionImageStore
}

func NewUpdateServiceRepo(
	serviceStore UpdateServiceStore,
	serviceVersionStore UpdateServiceVersionStore,
	imageRepo UpdateImageRepo,
	m2mVersionImages UpdateM2MVersionImageStore,
) *updateServiceRepo {
	return &updateServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
		imageRepo:           imageRepo,
		m2mVersionImages:    m2mVersionImages,
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
	existingService, err := repo.serviceStore.FindOne(
		ctx,
		map[string]interface{}{"id": serviceID},
		"Versions",
		"Images",
	)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	if existingService == nil {
		return nil, common.ErrEntityNotFound(models.ServiceEntityName, nil)
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

		existingServiceVersion, err := repo.serviceVersionStore.FindOne(
			ctx,
			map[string]interface{}{"id": serviceVersionId},
			"Images",
		)
		if err != nil {
			logger.AppLogger.Error(ctx, "version not found", zap.Error(err))
			return nil, err
		}

		serviceVersion := &models.ServiceVersion{
			SQLModel:      common.SQLModel{Id: serviceVersionId},
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

				img.ServiceVersions = append(img.ServiceVersions, serviceVersion)
				serviceVersion.Images = append(serviceVersion.Images, img)
			}
		}

		if input.ServiceVersion.VersionImages != nil {
			if len(*input.ServiceVersion.VersionImages) == 0 {
				serviceVersion.Images = []*models.Image{}
			} else {
				imageIDs := make([]uint32, len(*input.ServiceVersion.VersionImages))
				for _, versionImage := range *input.ServiceVersion.VersionImages {
					imageID, err := versionImage.GetLocalID(ctx)
					if err != nil {
						return nil, err
					}
					imageIDs = append(imageIDs, imageID)
				}

				images, err := repo.imageRepo.List(
					ctx,
					[]interface{}{imageIDs},
				)
				if err != nil {
					logger.AppLogger.Error(ctx, "version image not found", zap.Error(err))
					return nil, err
				}

				serviceVersion.Images = append(serviceVersion.Images, images...)
			}

		}

		if input.ServiceVersion.MainImageID != nil {
			mainImageID, err := input.ServiceVersion.GetMainImageLocalId(ctx)
			if err != nil {
				return nil, err
			}

			if _, hasThisImg := lo.Find(existingService.Images, func(img models.Image) bool {
				return img.Id == mainImageID
			}); !hasThisImg {
				return nil, common.ErrInvalidRequest(errors.New("main image not found"))
			}

			serviceVersion.MainImageID = mainImageID
		}

		// create new version as draft if the service and current version were published
		if existingService.Status == common.StatusActive &&
			existingServiceVersion.Status == common.StatusActive &&
			existingServiceVersion.PublishedDate != nil {
			_, hasDraft := lo.Find(existingService.Versions, func(version models.ServiceVersion) bool {
				return version.PublishedDate == nil
			})

			if hasDraft {
				return nil, serviceerror.ErrServiceDraftExisting(errors.New("service draft already existing"))
			}

			serviceVersion.Status = common.StatusInactive
			if err := repo.serviceVersionStore.CreateNewServiceVersion(ctx, serviceVersion); err != nil {
				return nil, common.ErrDB(err)
			}
		} else {
			if err := repo.serviceVersionStore.Update(ctx, serviceVersionId, serviceVersion); err != nil {
				return nil, common.ErrDB(err)
			}
		}

		existingService.ServiceVersionID = &serviceVersion.Id
		existingService.ServiceVersion = serviceVersion
	}

	return existingService, nil
}

func (repo *updateServiceRepo) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Service, error) {
	return repo.serviceStore.FindOne(ctx, conditions, moreInfo...)
}
