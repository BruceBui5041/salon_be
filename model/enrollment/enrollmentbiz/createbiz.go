package enrollmentbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/enrollment/enrollmentmodel"

	"github.com/shopspring/decimal"
)

type EnrollmentRepo interface {
	CreateNewEnrollment(
		ctx context.Context,
		userID, serviceID uint32,
		paymentID *uint32,
	) error
	CheckDuplicateEnrollment(ctx context.Context, userID, serviceID uint32) (bool, error)
	FindService(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.ServiceVersion, error)
}

type createEnrollmentBiz struct {
	enrollmentRepo EnrollmentRepo
}

func NewCreateEnrollmentBiz(enrollmentRepo EnrollmentRepo) *createEnrollmentBiz {
	return &createEnrollmentBiz{
		enrollmentRepo: enrollmentRepo,
	}
}

func (e *createEnrollmentBiz) CreateNewEnrollment(
	ctx context.Context,
	input *enrollmentmodel.CreateEnrollment,
) error {
	userUID, err := common.FromBase58(input.UserID)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	serviceUID, err := common.FromBase58(input.ServiceID)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	if err := e.validateInput(input); err != nil {
		return err
	}

	isDuplicate, err := e.enrollmentRepo.CheckDuplicateEnrollment(ctx, userUID.GetLocalID(), serviceUID.GetLocalID())
	if err != nil {
		return common.ErrCannotGetEntity(models.EnrollmentEntityName, err)
	}
	if isDuplicate {
		return errors.New("user is already enrolled in this service")
	}

	service, err := e.enrollmentRepo.FindService(ctx, map[string]interface{}{"id": serviceUID.GetLocalID()})
	if err != nil {
		return common.ErrCannotGetEntity(models.ServiceEntityName, err)
	}

	var paymentID *uint32
	if service.Price.Equal(decimal.Zero) {
		if input.PaymentID == "" {
			return errors.New("payment is required for non-free services")
		}
		paymentUID, err := common.FromBase58(input.PaymentID)
		if err != nil {
			return common.ErrInvalidRequest(err)
		}
		localPaymentID := paymentUID.GetLocalID()
		paymentID = &localPaymentID
	}

	err = e.enrollmentRepo.CreateNewEnrollment(ctx, userUID.GetLocalID(), serviceUID.GetLocalID(), paymentID)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	return nil
}

func (e *createEnrollmentBiz) validateInput(input *enrollmentmodel.CreateEnrollment) error {
	if input.UserID == "" {
		return errors.New("user ID is required")
	}

	if input.ServiceID == "" {
		return errors.New("service ID is required")
	}

	return nil
}
