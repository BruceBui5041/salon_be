package otprepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"time"
)

type VerifyOTPStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.OTP, error)
	Update(ctx context.Context, otp *models.OTP) error
}

type verifyRepo struct {
	store VerifyOTPStore
}

func NewVerifyRepo(store VerifyOTPStore) *verifyRepo {
	return &verifyRepo{store: store}
}

func (r *verifyRepo) FindOTPByUserID(ctx context.Context, userID uint32) (*models.OTP, error) {
	conditions := map[string]interface{}{
		"user_id":   userID,
		"passed_at": nil,
	}

	otp, err := r.store.FindOne(ctx, conditions)
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
	if err := r.store.Update(ctx, otp); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
