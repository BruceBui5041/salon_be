package notificationerror

import (
	"net/http"
	"salon_be/common"
)

const (
	errNotificationTitleEmpty   = "ErrNotificationTitleEmpty"
	errNotificationTitleTooLong = "ErrNotificationTitleTooLong"
	errNotificationContentEmpty = "ErrNotificationContentEmpty"
	errNotificationTypeEmpty    = "ErrNotificationTypeEmpty"
	errNotificationCannotCreate = "ErrNotificationCannotCreate"
)

func ErrNotificationTitleEmpty(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Notification title is required",
		err.Error(),
		errNotificationTitleEmpty,
	)
}

func ErrNotificationTitleTooLong(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Notification title must not exceed 255 characters",
		err.Error(),
		errNotificationTitleTooLong,
	)
}

func ErrNotificationContentEmpty(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Notification content is required",
		err.Error(),
		errNotificationContentEmpty,
	)
}

func ErrNotificationTypeEmpty(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Notification type is required",
		err.Error(),
		errNotificationTypeEmpty,
	)
}

func ErrNotificationCannotCreate(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Cannot create notification",
		err.Error(),
		errNotificationCannotCreate,
	)
}
