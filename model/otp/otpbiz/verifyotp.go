package otpbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/otp/otperror"
	"salon_be/model/otp/otpmodel"
)

type VerifyRepository interface {
	FindOTPByUserID(ctx context.Context, userID uint32) (*models.OTP, error)
	UpdateOTPStatus(ctx context.Context, otp *models.OTP) error
}

type verifyBiz struct {
	repo VerifyRepository
}

func NewVerifyBiz(repo VerifyRepository) *verifyBiz {
	return &verifyBiz{repo: repo}
}

func (biz *verifyBiz) VerifyOTP(ctx context.Context, input *otpmodel.VerifyOTPInput) error {
	otp, err := biz.repo.FindOTPByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if !otp.IsPassed(input.OTP) {
		return otperror.ErrOTPVerifyFailed(errors.New("verifying OTP failed"))
	}

	if err := biz.repo.UpdateOTPStatus(ctx, otp); err != nil {
		return common.ErrCannotUpdateEntity(models.OTPEntityName, err)
	}

	return nil
}
