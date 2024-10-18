package enrollmentutils

import (
	"context"
	"salon_be/component"
	"salon_be/component/cache"
	"salon_be/component/logger"
	"salon_be/watermill/messagemodel"

	"go.uber.org/zap"
)

func UpdateEnrollmentCache(
	ctx context.Context,
	appCtx component.AppContext,
	info *messagemodel.EnrollmentChangeInfo,
) error {
	enrollment := &cache.EnrollmentCache{
		UserId:            info.UserId,
		CourseId:          info.CourseId,
		CourseSlug:        info.CourseSlug,
		EnrollmentId:      info.EnrollmentId,
		PaymentId:         info.PaymentId,
		TransactionStatus: info.TransactionStatus,
	}

	appCache := appCtx.GetAppCache()

	if err := appCache.DeleteEnrollmentCache(ctx, info.CourseId, info.UserId); err != nil {
		logger.AppLogger.Error(ctx,
			"Failed to DeleteEnrollmentCache",
			zap.Error(err),
			zap.Any("updateUserCacheInfo", info),
		)
		return err
	}

	if err := appCache.SetEnrollmentCache(ctx, enrollment); err != nil {
		logger.AppLogger.Error(ctx,
			"Failed to SetEnrollmentCache",
			zap.Error(err),
			zap.Any("enrollment", enrollment),
		)
		return err
	}

	return nil
}
