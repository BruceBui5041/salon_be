package enrollmentrepo

import (
	"context"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/zap"
)

type CreateEnrollmentStore interface {
	Create(
		ctx context.Context,
		newEnrollment *models.Enrollment,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Enrollment, error)
}

type ServiceStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.ServiceVersion, error)
}

type createEnrollmentRepo struct {
	store        CreateEnrollmentStore
	serviceStore ServiceStore
	localPub     *gochannel.GoChannel
}

func NewCreateEnrollmentRepo(
	store CreateEnrollmentStore,
	serviceStore ServiceStore,
	localPub *gochannel.GoChannel,
) *createEnrollmentRepo {
	return &createEnrollmentRepo{
		store:        store,
		serviceStore: serviceStore,
		localPub:     localPub,
	}
}

func (repo *createEnrollmentRepo) CreateNewEnrollment(ctx context.Context, userID, serviceVersionID uint32, paymentID *uint32) error {
	newEnrollment := &models.Enrollment{
		UserID:           userID,
		ServiceVersionID: serviceVersionID,
		PaymentID:        paymentID,
	}

	enrollId, err := repo.store.Create(ctx, newEnrollment)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	enrollment, err := repo.store.FindOne(
		ctx,
		map[string]interface{}{"id": enrollId},
		"Service",
		"Payment",
	)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	tempUser := common.SQLModel{Id: userID}
	tempUser.GenUID(common.DbTypeUser)

	enrollment.Mask(false)
	enrollment.ServiceVersion.Mask(false)
	enrollment.Payment.Mask(false)

	// publish event to update user cache on local and dynamoDB
	updateCacheMsg := &messagemodel.EnrollmentChangeInfo{
		UserId:            tempUser.GetFakeId(),
		ServiceId:         enrollment.ServiceVersion.GetFakeId(),
		ServiceSlug:       enrollment.ServiceVersion.Slug,
		EnrollmentId:      enrollment.GetFakeId(),
		PaymentId:         enrollment.Payment.GetFakeId(),
		TransactionStatus: enrollment.Payment.TransactionStatus,
	}

	if err := watermill.PublishEnrollmentChange(ctx, repo.localPub, updateCacheMsg); err != nil {
		logger.AppLogger.Error(
			ctx,
			"cannot publish update user cache message",
			zap.Error(common.ErrInternal(err)),
			zap.Any("updateCacheMsg", updateCacheMsg),
		)
	}

	return nil
}

func (repo *createEnrollmentRepo) CheckDuplicateEnrollment(ctx context.Context, userID, serviceVersionID uint32) (bool, error) {
	enrollment, err := repo.store.FindOne(ctx, map[string]interface{}{
		"user_id":            userID,
		"service_version_id": serviceVersionID,
	})
	if err != nil {
		if err == common.RecordNotFound {
			return false, nil
		}
		return false, err
	}
	return enrollment != nil, nil
}

func (repo *createEnrollmentRepo) FindService(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...interface{},
) (*models.ServiceVersion, error) {
	return repo.serviceStore.FindOne(ctx, conditions, moreInfo...)
}
