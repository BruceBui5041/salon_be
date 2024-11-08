package otperror

import (
	"net/http"
	"salon_be/common"
)

const errOTPVerifyFailed = "ErrOTPVerifyFailed"

func ErrOTPVerifyFailed(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errOTPVerifyFailed,
	)
}
