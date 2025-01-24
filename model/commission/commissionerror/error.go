package commissionerror

import (
	"fmt"
	"net/http"
	"salon_be/common"
)

const (
	errCommissionPublished = "ErrCommissionPublished"
)

func ErrCommissionPublished() *common.AppError {
	err := fmt.Errorf("commission has been published and cannot be updated")
	return common.NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		err.Error(),
		err.Error(),
		errCommissionPublished,
	)
}
