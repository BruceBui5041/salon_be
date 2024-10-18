package enrollmentrepo

import (
	"context"
	"video_server/common"
	"video_server/component/logger"
	models "video_server/model"
	"video_server/watermill"
	"video_server/watermill/messagemodel"

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

type CourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type createEnrollmentRepo struct {
	store       CreateEnrollmentStore
	courseStore CourseStore
	localPub    *gochannel.GoChannel
}

func NewCreateEnrollmentRepo(
	store CreateEnrollmentStore,
	courseStore CourseStore,
	localPub *gochannel.GoChannel,
) *createEnrollmentRepo {
	return &createEnrollmentRepo{
		store:       store,
		courseStore: courseStore,
		localPub:    localPub,
	}
}

func (repo *createEnrollmentRepo) CreateNewEnrollment(ctx context.Context, userID, courseID uint32, paymentID *uint32) error {
	newEnrollment := &models.Enrollment{
		UserID:    userID,
		CourseID:  courseID,
		PaymentID: paymentID,
	}

	enrollId, err := repo.store.Create(ctx, newEnrollment)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	enrollment, err := repo.store.FindOne(
		ctx,
		map[string]interface{}{"id": enrollId},
		"Course",
		"Payment",
	)
	if err != nil {
		return common.ErrCannotCreateEntity(models.EnrollmentEntityName, err)
	}

	tempUser := common.SQLModel{Id: userID}
	tempUser.GenUID(common.DbTypeUser)

	enrollment.Mask(false)
	enrollment.Course.Mask(false)
	enrollment.Payment.Mask(false)

	// publish event to update user cache on local and dynamoDB
	updateCacheMsg := &messagemodel.EnrollmentChangeInfo{
		UserId:            tempUser.GetFakeId(),
		CourseId:          enrollment.Course.GetFakeId(),
		CourseSlug:        enrollment.Course.Slug,
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

func (repo *createEnrollmentRepo) CheckDuplicateEnrollment(ctx context.Context, userID, courseID uint32) (bool, error) {
	enrollment, err := repo.store.FindOne(ctx, map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
	})
	if err != nil {
		if err == common.RecordNotFound {
			return false, nil
		}
		return false, err
	}
	return enrollment != nil, nil
}

func (repo *createEnrollmentRepo) FindCourse(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...interface{},
) (*models.Course, error) {
	return repo.courseStore.FindOne(ctx, conditions, moreInfo...)
}
