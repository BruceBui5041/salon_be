package servicebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
)

type PublishServiceRepo interface {
	FindService(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Service, error)
	PublishService(ctx context.Context, serviceId uint32, versionId uint32) error
	FindServiceVersion(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.ServiceVersion, error)
}

type publishServiceBiz struct {
	repo PublishServiceRepo
}

func NewPublishServiceBiz(repo PublishServiceRepo) *publishServiceBiz {
	return &publishServiceBiz{repo: repo}
}

func (biz *publishServiceBiz) PublishService(
	ctx context.Context,
	requester common.Requester,
	data *servicemodel.PublishServiceRequest,
) error {
	serviceId, err := data.GetServiceLocalId()
	if err != nil {
		panic(err)
	}

	versionId, err := data.GetServiceVersionLocalId()
	if err != nil {
		panic(err)
	}

	// Get service with creator info
	service, err := biz.repo.FindService(
		ctx,
		map[string]interface{}{"id": serviceId},
	)
	if err != nil {
		return common.ErrCannotGetEntity(models.ServiceEntityName, err)
	}

	if service.ServiceVersionID != nil && *service.ServiceVersionID == versionId {
		return common.ErrInvalidRequest(errors.New("this version currently has been published"))
	}

	if service.CreatorID != requester.GetUserId() {
		return common.ErrNoPermission(errors.New("you are not the creator of this service"))
	}

	if err := biz.repo.PublishService(ctx, serviceId, versionId); err != nil {
		return common.ErrCannotUpdateEntity(models.ServiceEntityName, err)
	}

	return nil
}
