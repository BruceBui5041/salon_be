package servicerepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
)

type ServiceStore interface {
	CreateNewService(ctx context.Context, data *models.Service) error
	Update(ctx context.Context, serviceID uint32, data *models.Service) error
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
			ServiceID:     service.Id,
			Title:         input.ServiceVersion.Title,
			Description:   input.ServiceVersion.Description,
			CategoryID:    categoryUID.GetLocalID(),
			SubCategoryID: subCategoryUID.GetLocalID(),
			IntroVideoID:  introVideoUID,
			Thumbnail:     input.ServiceVersion.Thumbnail,
			Price:         input.ServiceVersion.Price.GetDecimal(),
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

		if err := repo.serviceStore.Update(ctx, service.Id, service); err != nil {
			return nil, common.ErrDB(err)
		}

		service.ServiceVersionID = &serviceVersion.Id
	}

	return service, nil
}
