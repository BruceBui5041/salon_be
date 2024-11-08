package otpmodel

type VerifyOTPInput struct {
	OTP    string `json:"otp" binding:"required"`
	UserID uint32 `json:"user_id"`
}
