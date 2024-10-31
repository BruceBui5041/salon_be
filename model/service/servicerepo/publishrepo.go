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

func (repo *publishServiceRepo) PublishService(
	ctx context.Context,
	serviceId uint32,
	versionId uint32,
) error {
	// Get service version to check status
	version, err := repo.serviceVersionStore.FindOne(
		ctx,
		map[string]interface{}{"id": versionId},
	)
	if err != nil {
		return err
	}

	// Set service and version status to active if inactive
	if version.Status == common.StatusInactive {
		// Update service status
		if err := repo.serviceStore.Update(
			ctx,
			serviceId,
			&models.Service{
				SQLModel: common.SQLModel{
					Status: common.StatusActive,
				},
			},
		); err != nil {
			return err
		}

		// Update version status and set published date
		now := time.Now()
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
