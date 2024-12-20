package userdevicebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/userdevice/userdevicemodel"
)

type UserDeviceRepo interface {
	CreateUserDevice(ctx context.Context, input *userdevicemodel.CreateUserDevice) (*models.UserDevice, error)
}

type createUserDeviceBiz struct {
	repo UserDeviceRepo
}

func NewCreateUserDeviceBiz(repo UserDeviceRepo) *createUserDeviceBiz {
	return &createUserDeviceBiz{repo: repo}
}

func (biz *createUserDeviceBiz) CreateUserDevice(ctx context.Context, input *userdevicemodel.CreateUserDevice) error {
	if input.FCMToken == "" {
		return common.ErrInvalidRequest(errors.New("fcm token is required"))
	}

	if len(input.FCMToken) > 250 {
		return common.ErrInvalidRequest(errors.New("fcm token must not exceed 250 characters"))
	}

	if input.Platform == "" {
		return common.ErrInvalidRequest(errors.New("platform is required"))
	}

	userDevice, err := biz.repo.CreateUserDevice(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.UserDeviceEntityName, err)
	}

	input.Id = userDevice.Id
	return nil
}
