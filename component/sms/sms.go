package sms

import (
	"context"
	"salon_be/component/sms/esms"
)

type ESMS interface {
	ESMSSendOTP(ctx context.Context, otpPayload *esms.OTPPayload) error
}

type smsClient struct {
	eSMS ESMS
}

type OTPMessage struct {
	Content     string
	PhoneNumber string
}

func NewSMSClient() *smsClient {
	return &smsClient{
		eSMS: esms.NewESMSClient(),
	}
}

func (client *smsClient) SendOTP(ctx context.Context, otpMessage OTPMessage) error {
	otpPayload := &esms.OTPPayload{
		Content: otpMessage.Content,
		Phone:   otpMessage.PhoneNumber,
	}

	return client.eSMS.ESMSSendOTP(ctx, otpPayload)
}
