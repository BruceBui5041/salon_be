package esms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"salon_be/common"
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
	IsUnicode   string `json:"IsUnicode"`
	Sandbox     string `json:"Sandbox"`
	CampaignID  string `json:"campaignid"`
	RequestID   string `json:"RequestId"`
	CallbackUrl string `json:"CallbackUrl"`
	SendDate    string `json:"SendDate"`
}

type ESMSResponse struct {
	CodeResult      string `json:"CodeResult"`
	SMSID           string `json:"SMSID"`
	CountRegenerate int    `json:"CountRegenerate"`
	ErrorMessage    string `json:"ErrorMessage"`
}

type eSMS struct {
	httpClient *http.Client
}

func NewESMSClient() *eSMS {
	return &eSMS{
		httpClient: &http.Client{},
	}
}

func (esms *eSMS) ESMSSendOTP(ctx context.Context, otpPayload *OTPPayload) (*ESMSResponse, error) {
	eSMSRestAPI := viper.GetString("OTP_SMS_REST_API")
	otpPayload.ApiKey = viper.GetString("OTP_SMS_API_KEY")
	otpPayload.SecretKey = viper.GetString("OTP_SMS_API_SECRET_KEY")
	otpPayload.Brandname = viper.GetString("OTP_SMS_BRANDNAME")
	otpPayload.SmsType = "2"
	otpPayload.IsUnicode = "0"
	otpPayload.Sandbox = "1"

	payloadJSON, err := json.Marshal(otpPayload)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to marshal payload", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest("POST", eSMSRestAPI, bytes.NewBuffer(payloadJSON))
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to create request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := esms.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to read response body", zap.Error(err))
		return nil, err
	}

	otpResponse := &ESMSResponse{}
	err = json.Unmarshal(body, &otpResponse)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to unmarshal response", zap.Error(err))
		return nil, err
	}

	if otpResponse.CodeResult != "100" {
		logger.AppLogger.Error(ctx, "Failed to send OTP", zap.Any("response", otpResponse))
		return nil, common.ErrInternal(errors.New(otpResponse.ErrorMessage))
	}

	logger.AppLogger.Info(ctx, "Send OTP successful", zap.Any("response", otpResponse))
	return otpResponse, nil
}
