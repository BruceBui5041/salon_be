// coupon/couponerror/error.go
package couponerror

import (
	"net/http"
	"salon_be/common"
)

const (
	errCouponInvalid          = "ErrCouponInvalid"
	errCouponExists           = "ErrCouponExists"
	errCouponExpired          = "ErrCouponExpired"
	errCouponUsageLimit       = "ErrCouponUsageLimit"
	errCouponHasUsaged        = "ErrCouponHasUsaged"
	errCouponHasBeenActivated = "ErrCouponHasBeenActivated"
)

func ErrCouponHasBeenActivated(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponHasBeenActivated,
	)
}

func ErrCouponHasUsaged(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponHasUsaged,
	)
}

func ErrCouponInvalid(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponInvalid,
	)
}

func ErrCouponExists(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponExists,
	)
}

func ErrCouponExpired(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponExpired,
	)
}

func ErrCouponUsageLimit(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCouponUsageLimit,
	)
}
