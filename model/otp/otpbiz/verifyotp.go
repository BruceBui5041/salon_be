package otpbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
	"salon_be/model/otp/otperror"
	"salon_be/model/otp/otpmodel"

	"github.com/jinzhu/copier"
)

type VerifyRepository interface {
	FindOTPByUserID(ctx context.Context, userID uint32) (*models.OTP, error)
	UpdateOTPStatus(ctx context.Context, otp *models.OTP) error
	FindUserByID(ctx context.Context, userID uint32) (*models.User, error)
}

type verifyBiz struct {
	repo          VerifyRepository
	tokenProvider tokenprovider.Provider
	hasher        hasher.Hasher
	expiry        int
}

func NewVerifyBiz(
	repo VerifyRepository,
	tokenProvider tokenprovider.Provider,
	hasher hasher.Hasher,
	expiry int,
) *verifyBiz {
	return &verifyBiz{
		repo:          repo,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

func (biz *verifyBiz) VerifyOTP(ctx context.Context, input *otpmodel.VerifyOTPInput) (*otpmodel.VerifyOTPResponse, error) {
	otp, err := biz.repo.FindOTPByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if !otp.IsPassed(input.OTP) {
		return nil, otperror.ErrOTPVerifyFailed(errors.New("verifying OTP failed"))
	}

	if err := biz.repo.UpdateOTPStatus(ctx, otp); err != nil {
		return nil, common.ErrCannotUpdateEntity(models.OTPEntityName, err)
	}

	// Get user details for token generation
	user, err := biz.repo.FindUserByID(ctx, input.UserID)
	if err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}

	// Generate token
	payload := tokenprovider.TokenPayload{
		UserId: int(user.Id),
		Roles:  user.Roles,
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	user.Mask(false)
	var userRes otpmodel.GetUserResponse
	copier.Copy(&userRes, user)

	return &otpmodel.VerifyOTPResponse{
		Token: accessToken,
		User:  userRes,
	}, nil
}
