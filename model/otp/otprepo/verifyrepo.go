package otprepo

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"time"

	"go.uber.org/zap"
)

type VerifyOTPStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.OTP, error)
	Update(ctx context.Context, otp *models.OTP) error
}

type VerifyUserStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.User, error)
	Update(ctx context.Context, userID uint32, data *models.User) error
}

type verifyRepo struct {
	optStore  VerifyOTPStore
	userStore VerifyUserStore
}

func NewVerifyRepo(
	optStore VerifyOTPStore,
	userStore VerifyUserStore,
) *verifyRepo {
	return &verifyRepo{optStore: optStore, userStore: userStore}
}

func (r *verifyRepo) FindOTPByUserID(ctx context.Context, userID uint32) (*models.OTP, error) {
	conditions := map[string]interface{}{
		"user_id":        userID,
		"passed_at":      nil,
		"expires_at > ?": time.Now().UTC(),
	}

	otp, err := r.optStore.FindOne(ctx, conditions)
	if err != nil {
		if err == common.RecordNotFound {
			return nil, common.ErrEntityNotFound(models.OTPEntityName, err)
		}
		return nil, common.ErrDB(err)
	}

	if otp.ExpiresAt.Before(time.Now().UTC()) {
		return nil, common.ErrInvalidRequest(errors.New("OTP expired"))
	}

	return otp, nil
}

func (r *verifyRepo) UpdateOTPStatus(ctx context.Context, otp *models.OTP) error {
	user, err := r.userStore.FindOne(ctx, map[string]interface{}{"id": otp.UserID})
	if err != nil {
		logger.AppLogger.Error(ctx, "get user error", zap.Error(err))
		return common.ErrEntityNotFound(models.UserEntityName, err)
	}

	if user.Status == common.StatusInactive {
		if err := r.userStore.Update(
			ctx,
			otp.UserID,
			&models.User{
				SQLModel: common.SQLModel{Status: common.StatusActive},
			},
		); err != nil {
			logger.AppLogger.Error(ctx, "update user status failed", zap.Error(err))
			return common.ErrDB(err)
		}
	}

	if err := r.optStore.Update(ctx, otp); err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (r *verifyRepo) FindUserByID(ctx context.Context, userID uint32) (*models.User, error) {
	return r.userStore.FindOne(ctx, map[string]interface{}{"id": userID}, "Roles", "UserProfile")
}
