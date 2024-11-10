package otpbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/component/sms"
	models "salon_be/model"
	"salon_be/model/otp/otperror"
	"salon_be/model/otp/otpmodel"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ResendOTPRepo interface {
	GetRecentOTPs(ctx context.Context, userId uint32, duration time.Duration) ([]*models.OTP, error)
	GetActiveOTP(ctx context.Context, userId uint32) (*models.OTP, error)
	GetUser(ctx context.Context, userId uint32) (*models.User, error)
	Create(ctx context.Context, data *models.OTP) error
	Update(ctx context.Context, data *models.OTP) error
}

type resendOTPBiz struct {
	repo      ResendOTPRepo
	smsClient component.SMSClient
}

func NewResendOTPBiz(
	repo ResendOTPRepo,
	smsClient component.SMSClient,
) *resendOTPBiz {
	return &resendOTPBiz{
		repo:      repo,
		smsClient: smsClient,
	}
}

func (biz *resendOTPBiz) ResendOTP(ctx context.Context, data *otpmodel.ResendOTPInput) error {
	activeOTP, err := biz.repo.GetActiveOTP(ctx, data.UserID)
	if err == nil && activeOTP != nil {
		return otperror.ErrActiveOTPExists(errors.New("an active OTP already exists"))
	}

	recentOTPs, err := biz.repo.GetRecentOTPs(ctx, data.UserID, time.Hour)
	if err != nil {
		return common.ErrDB(err)
	}

	if len(recentOTPs) >= 3 {
		return otperror.ErrOTPLimitExceeded(errors.New("exceeded maximum OTP attempts for this hour"))
	}

	user, err := biz.repo.GetUser(ctx, data.UserID)
	if err != nil {
		logger.AppLogger.Error(ctx, "get user error", zap.Error(err))
		return common.ErrDB(err)
	}

	newOTP := &models.OTP{UserID: data.UserID}

	newOTP.UUID = uuid.NewString()
	newOTP.ExpiresAt = time.Now().UTC().Add(5 * time.Minute)

	otp, err := generateOTP()
	if err != nil {
		return common.ErrInternal(err)
	}

	newOTP.OTP = otp

	if err := biz.repo.Create(ctx, newOTP); err != nil {
		logger.AppLogger.Error(ctx, "create OTP record failed", zap.Error(err))
		return common.ErrDB(err)
	}

	otpRes, err := biz.smsClient.SendOTP(ctx, sms.OTPMessage{
		Content:     "Your OTP is: " + otp,
		PhoneNumber: user.PhoneNumber,
	})
	if err != nil {
		logger.AppLogger.Error(ctx, "cannot send OTP", zap.Error(err))
		return common.ErrInternal(err)
	}

	newOTP.ESMSID = otpRes.SMSID
	if err := biz.repo.Update(ctx, newOTP); err != nil {
		logger.AppLogger.Error(ctx, "update OTP ESMSID failed", zap.Error(err))
		return err
	}

	return nil
}
