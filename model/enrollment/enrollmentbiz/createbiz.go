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
		userID, courseID uint32,
		paymentID *uint32,
	) error
	CheckDuplicateEnrollment(ctx context.Context, userID, courseID uint32) (bool, error)
	FindCourse(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
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

	courseUID, err := common.FromBase58(input.CourseID)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	if err := e.validateInput(input); err != nil {
		return err
	}

	isDuplicate, err := e.enrollmentRepo.CheckDuplicateEnrollment(ctx, userUID.GetLocalID(), courseUID.GetLocalID())
	if err != nil {
		return common.ErrCannotGetEntity(models.EnrollmentEntityName, err)
	}
	if isDuplicate {
		return errors.New("user is already enrolled in this course")
	}

	course, err := e.enrollmentRepo.FindCourse(ctx, map[string]interface{}{"id": courseUID.GetLocalID()})
	if err != nil {
		return common.ErrCannotGetEntity(models.CourseEntityName, err)
	}

	var paymentID *uint32
	if course.Price.Equal(decimal.Zero) {
		if input.PaymentID == "" {
			return errors.New("payment is required for non-free courses")
		}
		paymentUID, err := common.FromBase58(input.PaymentID)
		if err != nil {
			return common.ErrInvalidRequest(err)
		}
		localPaymentID := paymentUID.GetLocalID()
		paymentID = &localPaymentID
	}

	err = e.enrollmentRepo.CreateNewEnrollment(ctx, userUID.GetLocalID(), courseUID.GetLocalID(), paymentID)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	return nil
}

func (e *createEnrollmentBiz) validateInput(input *enrollmentmodel.CreateEnrollment) error {
	if input.UserID == "" {
		return errors.New("user ID is required")
	}

	if input.CourseID == "" {
		return errors.New("course ID is required")
	}

	return nil
}
