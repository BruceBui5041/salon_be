package servicerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"

	"github.com/shopspring/decimal"
)

type ServiceStore interface {
	CreateNewService(ctx context.Context, data *models.Service) error
}

type ServiceVersionStore interface {
	CreateNewServiceVersion(ctx context.Context, data *models.ServiceVersion) error
}

type createServiceRepo struct {
	serviceStore        ServiceStore
	serviceVersionStore ServiceVersionStore
}

func NewCreateServiceRepo(
	serviceStore ServiceStore,
	serviceVersionStore ServiceVersionStore,
) *createServiceRepo {
	return &createServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
	}
}

func (repo *createServiceRepo) CreateNewService(
	ctx context.Context,
	input *servicemodel.CreateService,
) (*models.Service, error) {
	service := &models.Service{
		CreatorID: input.CreatorID,
		Slug:      input.Slug,
	}

	if err := repo.serviceStore.CreateNewService(ctx, service); err != nil {
		return nil, common.ErrDB(err)
	}

	if input.ServiceVersion != nil {
		serviceVersion := &models.ServiceVersion{
			ServiceID:    service.Id,
			Title:        input.ServiceVersion.Title,
			Description:  input.ServiceVersion.Description,
			CategoryID:   input.ServiceVersion.CategoryID,
			IntroVideoID: input.ServiceVersion.IntroVideoID,
			Thumbnail:    input.ServiceVersion.Thumbnail,
			Price:        decimal.NewFromFloat(input.ServiceVersion.Price),
		}

		if input.ServiceVersion.DiscountedPrice != nil {
			discounted := decimal.NewFromFloat(*input.ServiceVersion.DiscountedPrice)
			serviceVersion.DiscountedPrice = &discounted
		}

		if err := repo.serviceVersionStore.CreateNewServiceVersion(ctx, serviceVersion); err != nil {
			return nil, common.ErrDB(err)
		}

		service.ServiceVersionID = &serviceVersion.Id
		service.ServiceVersion = serviceVersion
	}

	return service, nil
}
