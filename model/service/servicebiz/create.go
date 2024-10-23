package servicebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
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

		if input.ServiceVersion.CategoryID == 0 {
			return common.ErrInvalidRequest(errors.New("category ID is required"))
		}

		if input.ServiceVersion.Price <= 0 {
			return common.ErrInvalidRequest(errors.New("price must be greater than 0"))
		}
	}

	_, err := biz.repo.CreateNewService(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.ServiceEntityName, err)
	}

	return nil
}
