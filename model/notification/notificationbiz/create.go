// notification/notificationbiz/create.go
package notificationbiz

import (
	"context"
	"errors"
	models "salon_be/model"
	"salon_be/model/notification/notificationerror"
)

type NotificationRepo interface {
	CreateNotification(ctx context.Context, input *models.Notification) (*models.Notification, error)
}

type NotificationDetailRepo interface {
	CreateDetail(ctx context.Context, detail *models.NotificationDetail) (*models.NotificationDetail, error)
	FindDetails(ctx context.Context, conditions map[string]interface{}) ([]models.NotificationDetail, error)
}

type createNotificationBiz struct {
	notificationRepo       NotificationRepo
	notificationDetailRepo NotificationDetailRepo
}

func NewCreateNotificationBiz(
	repo NotificationRepo,
	detailRepo NotificationDetailRepo,
) *createNotificationBiz {
	return &createNotificationBiz{
		notificationRepo:       repo,
		notificationDetailRepo: detailRepo,
	}
}

func (biz *createNotificationBiz) CreateNotification(
	ctx context.Context,
	input *models.Notification,
) error {
	if input.Type == "" {
		return notificationerror.ErrNotificationTypeEmpty(
			errors.New("type is required"),
		)
	}

	notification, err := biz.notificationRepo.CreateNotification(ctx, input)
	if err != nil {
		return notificationerror.ErrNotificationCannotCreate(
			errors.New(err.Error()),
		)
	}

	input.Id = notification.Id
	return nil
}
