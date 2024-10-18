package paymentbiz

import (
	"context"
	"errors"
	"video_server/common"
	models "video_server/model"
	"video_server/model/payment/paymentmodel"

	"github.com/jinzhu/copier"
)

type PaymentRepo interface {
	CreateNewPayment(
		ctx context.Context,
		input *paymentmodel.CreatePayment,
	) (*models.Payment, error)
	CheckDuplicatePayment(ctx context.Context, userID uint32, transactionID string) (bool, error)
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Payment, error)
}

type EnrollmentRepo interface {
	CreateNewEnrollment(ctx context.Context, userID, courseID uint32, paymentID *uint32) error
	CheckDuplicateEnrollment(ctx context.Context, userID, courseID uint32) (bool, error)
}

type createPaymentBiz struct {
	paymentRepo    PaymentRepo
	enrollmentRepo EnrollmentRepo
}

func NewCreatePaymentBiz(paymentRepo PaymentRepo, enrollmentRepo EnrollmentRepo) *createPaymentBiz {
	return &createPaymentBiz{
		paymentRepo:    paymentRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (biz *createPaymentBiz) CreateNewPayment(
	ctx context.Context,
	input *paymentmodel.CreatePayment,
) (*paymentmodel.PaymentResponse, error) {
	if err := biz.validateInput(input); err != nil {
		return nil, err
	}

	duplicatePayment, err := biz.paymentRepo.CheckDuplicatePayment(ctx, input.UserID, input.TransactionID)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.PaymentEntityName, err)
	}
	if duplicatePayment {
		return nil, errors.New("duplicate payment found")
	}

	payment, err := biz.paymentRepo.CreateNewPayment(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.PaymentEntityName, err)
	}

	if payment.TransactionStatus == "completed" {
		for _, courseID := range input.CourseIDs {
			courseUID, err := common.FromBase58(courseID)
			if err != nil {
				return nil, common.ErrInvalidRequest(err)
			}
			duplicateEnrollment, err := biz.enrollmentRepo.CheckDuplicateEnrollment(ctx, input.UserID, courseUID.GetLocalID())
			if err != nil {
				return nil, common.ErrEntityExisted(models.EnrollmentEntityName, err)
			}
			if !duplicateEnrollment {
				err = biz.enrollmentRepo.CreateNewEnrollment(ctx, payment.UserID, courseUID.GetLocalID(), &payment.Id)
				if err != nil {
					return nil, common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
				}
			}
		}
	}

	payment, err = biz.paymentRepo.FindOne(
		ctx,
		map[string]interface{}{"id": payment.Id},
		"Enrollments.Course.Creator",
		"Enrollments.Course.Category",
	)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.PaymentEntityName, err)
	}

	payment.Mask(false)

	var res paymentmodel.PaymentResponse
	copier.Copy(&res, payment)
	return &res, nil
}

func (biz *createPaymentBiz) validateInput(input *paymentmodel.CreatePayment) error {
	if input.UserID == 0 {
		return errors.New("user ID is required")
	}
	if input.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if input.Currency == "" {
		return errors.New("currency is required")
	}
	if input.PaymentMethod == "" {
		return errors.New("payment method is required")
	}
	if len(input.CourseIDs) == 0 {
		return errors.New("at least one course ID is required")
	}
	return nil
}
