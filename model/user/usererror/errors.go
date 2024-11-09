package usererror

import (
	"errors"
	"net/http"
	"salon_be/common"
)

const errUserDraftExisting = "ErrUserMissingRequiredField"

var (
	ErrUserPhoneNumberNotFound = common.NewCustomError(
		errors.New("phone number is invalid"),
		"phone number not found",
		"ErrUserPhoneNumberNotFound",
	)

	ErrUsernameOrPasswordInvalid = common.NewCustomError(
		errors.New("username or password is invalid"),
		"username or password is invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailIsAlreadyExisted = common.NewCustomError(
		errors.New("email is alread existed"),
		"email is alread existed",
		"ErrEmailIsAlreadyExisted",
	)
)

func ErrUserMissionRequireField(err error) *common.AppError {
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errUserDraftExisting,
	)
}
