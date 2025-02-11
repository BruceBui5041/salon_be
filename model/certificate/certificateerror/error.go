package certificateerror

import (
	"net/http"
	"salon_be/common"
)

const (
	errInvalidFileType = "ErrInvalidFileType"
	errFileTooLarge    = "ErrFileTooLarge"
)

func ErrInvalidFileType(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"File must be in PDF format",
		"File must be in PDF format",
		errInvalidFileType,
	)
}

func ErrFileTooLarge(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"File size must be less than 5MB",
		"File size must be less than 5MB",
		errFileTooLarge,
	)
}
