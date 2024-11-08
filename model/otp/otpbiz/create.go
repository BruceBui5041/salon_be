package otpbiz

import (
	"context"
	"crypto/rand"
	"math/big"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/otp/otpmodel"
	"time"
)

type CreateOTPStore interface {
	Create(ctx context.Context, data *models.OTP) error
}

type createOTPBiz struct {
	store CreateOTPStore
}

func NewCreateOTPBiz(store CreateOTPStore) *createOTPBiz {
	return &createOTPBiz{store: store}
}

func (biz *createOTPBiz) CreateOTP(ctx context.Context, data *otpmodel.CreateOTPInput) error {
	if data == nil {
		return nil
	}

	newOTP := &models.OTP{
		UserID: data.UserID,
	}

	newOTP.ExpiresAt = time.Now().Add(5 * time.Minute)

	otp, err := generateOTP()
	if err != nil {
		return common.ErrInternal(err)
	}

	newOTP.OTP = otp

	if err := biz.store.Create(ctx, newOTP); err != nil {
		return err
	}

	return nil
}

func generateOTP() (string, error) {
	const otpLength = 6
	const digits = "0123456789"

	otp := make([]byte, otpLength)

	for i := 0; i < otpLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			// In case of error, return a default OTP (you might want to handle this differently)
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}

	return string(otp), nil
}
