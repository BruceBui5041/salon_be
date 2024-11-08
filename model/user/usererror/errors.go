package usererror

import (
	"net/http"
	"salon_be/common"
)

const errUserDraftExisting = "ErrUserMissingRequiredField"

func ErrUserMissionRequireField(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errUserDraftExisting,
	)
}
