package otpmodel

import (
	"salon_be/common"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
)

type VerifyOTPInput struct {
	OTP    string `json:"otp" binding:"required"`
	UserID uint32 `json:"user_id"`
}

type GetUserResponse struct {
	common.SQLModel `json:",inline"`
	Email           string              `json:"email"`
	PhoneNumber     string              `json:"phone_number"`
	UserProfile     *models.UserProfile `json:"user_profile"`
	Roles           []models.Role       `json:"roles"`
	Status          int                 `json:"status"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

type VerifyOTPResponse struct {
	Token     *tokenprovider.Token `json:"token"`
	User      GetUserResponse      `json:"user"`
	Challenge string               `json:"challenge,omitempty"`
}
