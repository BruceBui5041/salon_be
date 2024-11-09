package otpmodel

type ResendOTPInput struct {
	UserID uint32 `json:"user_id" binding:"required"`
}
