package otpmodel

type CreateOTPInput struct {
	UserID uint32 `json:"user_id" binding:"required"`
}
