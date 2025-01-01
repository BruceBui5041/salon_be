package userdevicerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/userdevice/userdevicemodel"
)

type CreateUserDeviceStore interface {
	Create(ctx context.Context, data *models.UserDevice) (*models.UserDevice, error)
	Update(
		ctx context.Context,
		id uint32,
		data *models.UserDevice,
	) error
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.UserDevice, error)
}

type createUserDeviceRepo struct {
	store CreateUserDeviceStore
}

func NewCreateUserDeviceRepo(store CreateUserDeviceStore) *createUserDeviceRepo {
	return &createUserDeviceRepo{store: store}
}

func (repo *createUserDeviceRepo) CreateUserDevice(ctx context.Context, input *userdevicemodel.CreateUserDevice) (*models.UserDevice, error) {
	// Check if device already exists for this user
	existing, _ := repo.store.FindOne(ctx, map[string]interface{}{
		"user_id":  input.UserID,
		"platform": input.Platform,
	})

	if existing != nil {
		if existing.FCMToken != input.FCMToken {
			existing.FCMToken = input.FCMToken
			err := repo.store.Update(ctx, existing.Id, existing)
			if err != nil {
				return nil, common.ErrDB(err)
			}

			return repo.store.FindOne(ctx, map[string]interface{}{
				"id": existing.Id,
			})

		}
		return existing, nil // Return existing device if found
	}

	userDevice := &models.UserDevice{
		FCMToken: input.FCMToken,
		Platform: input.Platform,
		UserID:   input.UserID,
	}

	_, err := repo.store.Create(ctx, userDevice)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	newUserDevice, err := repo.store.Create(ctx, userDevice)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return newUserDevice, nil
}
