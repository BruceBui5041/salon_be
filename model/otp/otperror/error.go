package otperror

import (
	"net/http"
	"salon_be/common"
)

const (
	errOTPVerifyFailed  = "ErrOTPVerifyFailed"
	errActiveOTPExists  = "ErrActiveOTPExists"
	errOTPLimitExceeded = "ErrOTPLimitExceeded"
)

func ErrOTPVerifyFailed(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errOTPVerifyFailed,
	)
}

func ErrActiveOTPExists(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errActiveOTPExists,
	)
}

func ErrOTPLimitExceeded(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusTooManyRequests,
		err,
		err.Error(),
		err.Error(),
		errOTPLimitExceeded,
	)
}
