package esms

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"salon_be/component/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type OTPPayload struct {
	ApiKey      string `json:"ApiKey"`
	Content     string `json:"Content"`
	Phone       string `json:"Phone"`
	SecretKey   string `json:"SecretKey"`
	Brandname   string `json:"Brandname"`
	SmsType     string `json:"SmsType"`
	IsUnicode   bool   `json:"IsUnicode"`
	Sandbox     bool   `json:"Sandbox"`
	CampaignID  string `json:"campaignid"`
	RequestID   string `json:"RequestId"`
	CallbackUrl string `json:"CallbackUrl"`
	SendDate    int64  `json:"SendDate"`
}

type eSMS struct {
	httpClient *http.Client
}

func NewESMSClient() *eSMS {
	return &eSMS{
		httpClient: &http.Client{},
	}
}

func (esms *eSMS) ESMSSendOTP(ctx context.Context, otpPayload *OTPPayload) error {
	eSMSRestAPI := viper.GetString("OTP_SMS_REST_API")
	otpPayload.ApiKey = viper.GetString("OTP_SMS_API_KEY")
	otpPayload.SecretKey = viper.GetString("OTP_SMS_SECRET_KEY")
	otpPayload.Brandname = viper.GetString("OTP_SMS_BRANDNAME")

	payloadJSON, err := json.Marshal(otpPayload)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to marshal payload", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", eSMSRestAPI, bytes.NewBuffer(payloadJSON))
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to create request", zap.Error(err))
		return err
	}

	resp, err := esms.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
