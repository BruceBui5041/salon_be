package servicerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"time"
)

type PublishServiceStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Service, error)
	Update(ctx context.Context, serviceID uint32, updates *models.Service) error
}

type PublishServiceVersionStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.ServiceVersion, error)
	Update(ctx context.Context, versionID uint32, updates *models.ServiceVersion) error
}

type publishServiceRepo struct {
	serviceStore        PublishServiceStore
	serviceVersionStore PublishServiceVersionStore
}

func NewPublishServiceRepo(
	serviceStore PublishServiceStore,
	serviceVersionStore PublishServiceVersionStore,
) *publishServiceRepo {
	return &publishServiceRepo{
		serviceStore:        serviceStore,
		serviceVersionStore: serviceVersionStore,
	}
}

func (repo *publishServiceRepo) FindService(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Service, error) {
	return repo.serviceStore.FindOne(ctx, conditions, moreInfo...)
}

func (repo *publishServiceRepo) FindServiceVersion(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.ServiceVersion, error) {
	return repo.serviceVersionStore.FindOne(ctx, conditions, moreInfo...)
}

func (repo *publishServiceRepo) PublishService(
	ctx context.Context,
	serviceId uint32,
	versionId uint32,
) error {
	service, err := repo.serviceStore.FindOne(
		ctx,
		map[string]interface{}{"id": serviceId},
	)
	if err != nil {
		return err
	}

	version, err := repo.serviceVersionStore.FindOne(
		ctx,
		map[string]interface{}{"id": versionId},
	)
	if err != nil {
		return err
	}

	if err := repo.serviceStore.Update(
		ctx,
		serviceId,
		&models.Service{
			SQLModel: common.SQLModel{
				Status: common.StatusActive,
			},
			ServiceVersionID: &versionId,
		},
	); err != nil {
		return err
	}

	if err := repo.serviceVersionStore.Update(ctx, versionId, &models.ServiceVersion{
		SQLModel: common.SQLModel{
			Status: common.StatusActive,
		},
	}); err != nil {
		return err
	}

	if service.ServiceVersionID != nil {
		if err := repo.serviceVersionStore.Update(
			ctx,
			*service.ServiceVersionID,
			&models.ServiceVersion{
				SQLModel: common.SQLModel{
					Status: common.StatusInactive,
				},
			},
		); err != nil {
			return err
		}
	}

	if version.PublishedDate == nil {
		now := time.Now().UTC()
		if err := repo.serviceVersionStore.Update(
			ctx,
			versionId,
			&models.ServiceVersion{
				SQLModel: common.SQLModel{
					Status: common.StatusActive,
				},
				PublishedDate: &now,
			},
		); err != nil {
			return err
		}
	}

	return nil
}
