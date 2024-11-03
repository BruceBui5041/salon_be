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
