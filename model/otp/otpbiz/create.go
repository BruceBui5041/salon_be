package otpbiz

import (
	"context"
	"crypto/rand"
	"math/big"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/component/sms"
	models "salon_be/model"
	"salon_be/model/otp/otpmodel"
	"time"

	"go.uber.org/zap"
)

type CreateOTPUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type CreateOTPStore interface {
	Create(ctx context.Context, data *models.OTP) error
}

type createOTPBiz struct {
	optStore  CreateOTPStore
	userStore CreateOTPUserStore
	smsClient component.SMSClient
}

func NewCreateOTPBiz(
	optStore CreateOTPStore,
	userStore CreateOTPUserStore,
	smsClient component.SMSClient,
) *createOTPBiz {
	return &createOTPBiz{
		optStore:  optStore,
		smsClient: smsClient,
		userStore: userStore,
	}
}

func (biz *createOTPBiz) CreateOTP(ctx context.Context, data *otpmodel.CreateOTPInput) error {
	if data == nil {
		return nil
	}

	newOTP := &models.OTP{
		UserID: data.UserID,
	}

	newOTP.ExpiresAt = time.Now().UTC().Add(5 * time.Minute)

	otp, err := generateOTP()
	if err != nil {
		return common.ErrInternal(err)
	}

	newOTP.OTP = otp

	if err := biz.optStore.Create(ctx, newOTP); err != nil {
		logger.AppLogger.Error(ctx, "create OTP record failed", zap.Error(err))
		return err
	}

	user, err := biz.userStore.FindOne(ctx, map[string]interface{}{"id": data.UserID})
	if err != nil {
		logger.AppLogger.Error(ctx, "get user error", zap.Error(err))
		return common.ErrDB(err)
	}

	err = biz.smsClient.SendOTP(ctx, sms.OTPMessage{
		Content:     "Your OTP is: " + otp,
		PhoneNumber: user.PhoneNumber,
	})
	if err != nil {
		logger.AppLogger.Error(ctx, "cannot send OTP", zap.Error(err))
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
