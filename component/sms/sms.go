package sms

import (
	"context"
	"salon_be/component/sms/esms"
)

type ESMS interface {
	ESMSSendOTP(ctx context.Context, otpPayload *esms.OTPPayload) (*esms.ESMSResponse, error)
}

type smsClient struct {
	eSMS ESMS
}

type OTPMessage struct {
	UUID        string
	Content     string
	PhoneNumber string
}

func NewSMSClient() *smsClient {
	return &smsClient{
		eSMS: esms.NewESMSClient(),
	}
}

func (client *smsClient) SendOTP(ctx context.Context, otpMessage OTPMessage) (*esms.ESMSResponse, error) {
	otpPayload := &esms.OTPPayload{
		RequestID: otpMessage.UUID,
		Content:   otpMessage.Content,
		Phone:     otpMessage.PhoneNumber,
	}

	return client.eSMS.ESMSSendOTP(ctx, otpPayload)
}
