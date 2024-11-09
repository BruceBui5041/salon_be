package otprepo

import (
	"context"
	models "salon_be/model"
	"time"
)

type ResendOTPStore interface {
	List(
		ctx context.Context,
		conditions []interface{},
		moreKeys ...string,
	) ([]*models.OTP, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.OTP, error)
	Create(ctx context.Context, data *models.OTP) error
	Update(ctx context.Context, updates *models.OTP) error
}

type UserStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.User, error)
}

type resendRepo struct {
	store     ResendOTPStore
	userStore UserStore
}

func NewResendRepo(store ResendOTPStore, userStore UserStore) *resendRepo {
	return &resendRepo{
		store:     store,
		userStore: userStore,
	}
}

func (repo *resendRepo) GetRecentOTPs(
	ctx context.Context,
	userId uint32,
	duration time.Duration,
) ([]*models.OTP, error) {
	fromTime := time.Now().UTC().Add(-duration)

	conditions := []interface{}{
		"user_id = ? AND created_at >= ?",
		userId,
		fromTime,
	}

	return repo.store.List(ctx, conditions)
}

func (repo *resendRepo) GetActiveOTP(
	ctx context.Context,
	userId uint32,
) (*models.OTP, error) {
	currentTime := time.Now().UTC()
	conditions := map[string]interface{}{
		"user_id":        userId,
		"passed_at":      nil,
		"expires_at > ?": currentTime,
	}

	return repo.store.FindOne(ctx, conditions)
}

func (repo *resendRepo) GetUser(
	ctx context.Context,
	userId uint32,
) (*models.User, error) {
	conditions := map[string]interface{}{
		"id": userId,
	}

	return repo.userStore.FindOne(ctx, conditions)
}

func (repo *resendRepo) Create(ctx context.Context, data *models.OTP) error {
	return repo.store.Create(ctx, data)
}

func (repo *resendRepo) Update(ctx context.Context, data *models.OTP) error {
	return repo.store.Update(ctx, data)
}
