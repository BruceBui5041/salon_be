package servicebiz

import (
	"context"
	"errors"
	"fmt"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"

	"github.com/shopspring/decimal"
)

type ServiceRepo interface {
	CreateNewService(ctx context.Context, input *servicemodel.CreateService) (*models.Service, error)
}

type createServiceBiz struct {
	repo ServiceRepo
}

func NewCreateServiceBiz(repo ServiceRepo) *createServiceBiz {
	return &createServiceBiz{repo: repo}
}

func (biz *createServiceBiz) CreateNewService(ctx context.Context, input *servicemodel.CreateService) error {
	if input.CreatorID == 0 {
		return common.ErrInvalidRequest(errors.New("creator ID is required"))
	}

	if input.Slug == "" {
		return common.ErrInvalidRequest(errors.New("slug is required"))
	}

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

		price := input.ServiceVersion.Price.GetDecimal()
		if price.Equal(decimal.Zero) {
			return common.ErrInvalidRequest(errors.New("price must be greater than 0"))
		}

		if input.ServiceVersion.DiscountedPrice != nil &&
			input.ServiceVersion.DiscountedPrice.Decimal.GreaterThanOrEqual(price) {
			return common.ErrInvalidRequest(fmt.Errorf("discount price must be less than price. Price %s", price.String()))
		}

		if input.ServiceVersion.Duration < 15 {
			return common.ErrInvalidRequest(errors.New("duration must be at least 15 minutes"))
		}

		requester, ok := ctx.Value(common.CurrentUser).(common.Requester)
		if !ok {
			return common.ErrInvalidRequest(errors.New("requester not found"))
		}

		if len(input.ServiceVersion.ServiceMenIds) != 0 &&
			!requester.IsAdmin() &&
			!requester.IsGroupProviderAdmin() {
			return common.ErrInvalidRequest(
				errors.New("only admin and group provider can update service men"),
			)
		}

		if input.ServiceVersion.OwnerID != nil && !requester.IsAdmin() {
			return common.ErrInvalidRequest(
				errors.New("only admin can assign owner"),
			)
		}

		if input.ServiceVersion.GroupProviderID != nil &&
			!requester.IsAdmin() {
			return common.ErrInvalidRequest(
				errors.New("only admin can assign group provider"),
			)
		}
	}

	newService, err := biz.repo.CreateNewService(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.ServiceEntityName, err)
	}

	if newService.ServiceVersionID == nil {
		return common.ErrInvalidRequest(errors.New("service must have service version"))
	}

	return nil
}
