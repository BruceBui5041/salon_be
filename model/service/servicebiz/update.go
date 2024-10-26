package servicebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
	"salon_be/utils/customtypes"

	"github.com/shopspring/decimal"
)

type UpdateServiceRepo interface {
	UpdateService(ctx context.Context, input *servicemodel.UpdateService) (*models.Service, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Service, error)
}

type updateServiceBiz struct {
	repo UpdateServiceRepo
}

func NewUpdateServiceBiz(repo UpdateServiceRepo) *updateServiceBiz {
	return &updateServiceBiz{repo: repo}
}

func (biz *updateServiceBiz) UpdateService(ctx context.Context, input *servicemodel.UpdateService) error {
	if input.ServiceVersion != nil {
		if input.ServiceVersion.Title == "" {
			return common.ErrInvalidRequest(errors.New("service version title is required"))
		}

		if input.ServiceVersion.CategoryID == "" {
			return common.ErrInvalidRequest(errors.New("category ID is required"))
		}

		if input.ServiceVersion.SubCategoryID == "" {
			return common.ErrInvalidRequest(errors.New("sub category ID is required"))
		}

		if input.ServiceVersion.Price != customtypes.DecimalString(decimal.Zero) {
			price := input.ServiceVersion.Price.GetDecimal()
			if price.Equal(decimal.Zero) {
				return common.ErrInvalidRequest(errors.New("price must be greater than 0"))
			}
		}

		if input.ServiceVersion.DiscountedPrice != nil {
			service, err := biz.repo.FindOne(
				ctx,
				map[string]interface{}{"id": input.ServiceID},
				"ServiceVersion",
			)
			if err != nil {
				return common.ErrDB(err)
			}
			if input.ServiceVersion.DiscountedPrice.Decimal.GreaterThanOrEqual(service.ServiceVersion.Price) {
				return common.ErrInvalidRequest(errors.New("discount price must be less than price"))
			}
		}

		if input.ServiceVersion.Duration < uint32(900) {
			return common.ErrInvalidRequest(errors.New("duration must be at least 15 minutes"))
		}
	}

	updatedService, err := biz.repo.UpdateService(ctx, input)
	if err != nil {
		return common.ErrCannotUpdateEntity(models.ServiceEntityName, err)
	}

	if updatedService.ServiceVersionID == nil {
		return common.ErrInvalidRequest(errors.New("service must have service version"))
	}

	return nil
}
