package notificationrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type NotificationStore interface {
	Create(ctx context.Context, data *models.Notification) (*models.Notification, error)
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Notification, error)
	Update(ctx context.Context, id uint32, data *models.Notification) error
}

type NotificationDetailStore interface {
	Create(ctx context.Context, data *models.NotificationDetail) (*models.NotificationDetail, error)
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.NotificationDetail, error)
	Update(ctx context.Context, conditions map[string]interface{}, data *models.NotificationDetail) error
}

type createNotificationRepo struct {
	store       NotificationStore
	detailStore NotificationDetailStore
}

func NewCreateNotificationRepo(
	store NotificationStore,
	detailStore NotificationDetailStore,
) *createNotificationRepo {
	return &createNotificationRepo{
		store:       store,
		detailStore: detailStore,
	}
}

func (repo *createNotificationRepo) CreateNotification(
	ctx context.Context,
	input *models.Notification,
) (*models.Notification, error) {
	notification := &models.Notification{
		Title:     input.Title,
		Content:   input.Content,
		Type:      input.Type,
		Scheduled: input.Scheduled,
		Metadata:  input.Metadata,
		BookingID: input.BookingID,
	}

	newNotification, err := repo.store.Create(ctx, notification)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return newNotification, nil
}

func (repo *createNotificationRepo) CreateDetail(
	ctx context.Context,
	detail *models.NotificationDetail,
) (*models.NotificationDetail, error) {
	newDetail, err := repo.detailStore.Create(ctx, detail)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return newDetail, nil
}

func (repo *createNotificationRepo) FindDetails(
	ctx context.Context,
	conditions map[string]interface{},
) ([]models.NotificationDetail, error) {
	details, err := repo.detailStore.Find(ctx, conditions)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return details, nil
}
