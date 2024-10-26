package servicerepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
	"salon_be/utils/customtypes"

	"github.com/shopspring/decimal"
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

type updateServiceRepo struct {
	serviceStore        UpdateServiceStore
	serviceVersionStore UpdateServiceVersionStore
}

func NewUpdateServiceRepo(
	serviceStore UpdateServiceStore,
	serviceVersionStore UpdateServiceVersionStore,
) *updateServiceRepo {
	return &updateServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
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

	service := &models.Service{
		ServiceVersionID: &serviceVersionId,
	}

	if err := repo.serviceStore.Update(ctx, serviceID, service); err != nil {
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

		var introVideoUID *uint32
		if input.ServiceVersion.IntroVideoID != "" {
			uid, err := common.DecomposeUID(input.ServiceVersion.IntroVideoID)
			if err != nil {
				return nil, common.ErrInvalidRequest(errors.New("invalid intro video ID"))
			}
			localId := uid.GetLocalID()
			introVideoUID = &localId
		}

		serviceVersion := &models.ServiceVersion{
			ServiceID:     serviceID,
			Title:         input.ServiceVersion.Title,
			Description:   input.ServiceVersion.Description,
			CategoryID:    categoryUID.GetLocalID(),
			SubCategoryID: subCategoryUID.GetLocalID(),
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
