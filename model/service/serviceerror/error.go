package serviceerror

import (
	"net/http"
	"salon_be/common"
)

const errServiceDraftExisting = "ErrServiceDraftExisting"

func ErrServiceDraftExisting(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Service draft already existing",
		err.Error(),
		errServiceDraftExisting,
	)
}

func ErrServiceVersionAlreadyPublished(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"This service version already published",
		err.Error(),
		errServiceDraftExisting,
	)
}
