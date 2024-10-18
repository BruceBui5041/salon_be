package middleware

import (
	"context"
	"errors"
	"video_server/common"
	"video_server/component"
	"video_server/component/cache"
	"video_server/component/logger"
	models "video_server/model"
	"video_server/model/course/coursestore"
	"video_server/model/enrollment/enrollmentstore"
	"video_server/model/video/videostore"
	"video_server/watermill"
	"video_server/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ErrPaymentNotCompleted(err error) *common.AppError {
	return common.NewCustomError(err, "payment was not completed", "ErrPaymentNotCompleted")
}

func EnrollmentCheck(appCtx component.AppContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		courseSlug := c.Param("course_slug")
		if courseSlug == "" {
			courseSlug = c.Query("course_slug")
		}

		videoIdFromCtx := c.Param("video_id")
		if videoIdFromCtx == "" {
			videoIdFromCtx = c.Query("video_id")
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}
		requester.Mask(false)

		// Try to get enrollment from cache first
		appCache := appCtx.GetAppCache()
		enrollmentCache, err := appCache.GetEnrollmentCache(c.Request.Context(), courseSlug, requester.GetFakeId())

		if err != nil {
			logger.AppLogger.Error(c, "Failed to get enrollment from cache", zap.Error(err))
			panic(common.ErrEntityNotFound(models.EnrollmentEntityName, err))
		}

		if enrollmentCache != nil {
			// Enrollment found in cache
			if enrollmentCache.TransactionStatus != "completed" {
				panic(ErrPaymentNotCompleted(errors.New("payment was not completed")))
			}
			c.Next()
			return
		}

		// If not in cache, fallback to database
		db := appCtx.GetMainDBConnection()
		enrollmentStore := enrollmentstore.NewSQLStore(db)
		courseStore := coursestore.NewSQLStore(db)

		course, err := courseStore.FindOne(c.Request.Context(), map[string]interface{}{"slug": courseSlug})
		if err != nil {
			panic(common.ErrDB(err))
		}
		course.Mask(false)

		videoCacheInfo, err := getVideoCacheInfo(c.Request.Context(), db, appCache, courseSlug, videoIdFromCtx)
		if err != nil {
			panic(common.ErrInternal(err))
		}

		// this mean the request video is the author
		if videoCacheInfo.CourseSlug == courseSlug && course.CreatorID == requester.GetUserId() {
			c.Next()
			return
		}

		enrollment, err := enrollmentStore.FindOne(
			c,
			map[string]interface{}{
				"user_id":   requester.GetUserId(),
				"course_id": course.Id,
			},
			"Payment",
		)

		if err != nil {
			panic(common.ErrDB(err))
		}

		enrollment.Mask(false)
		if enrollment.Payment != nil {
			enrollment.Payment.Mask(false)
		}

		if enrollment == nil {
			panic(common.ErrNoPermission(errors.New("not found your enrollment for this course")))
		}

		if enrollment.Payment.TransactionStatus != "completed" {
			panic(ErrPaymentNotCompleted(errors.New("payment was not completed")))
		}

		logger.AppLogger.Info(c.Request.Context(), "enrollment", zap.Any("enrollment", enrollment))

		// Publish enrollment change message
		updateCacheMsg := &messagemodel.EnrollmentChangeInfo{
			UserId:            requester.GetFakeId(),
			CourseId:          course.GetFakeId(),
			CourseSlug:        courseSlug,
			EnrollmentId:      enrollment.GetFakeId(),
			PaymentId:         enrollment.Payment.GetFakeId(),
			TransactionStatus: enrollment.Payment.TransactionStatus,
		}

		if err := watermill.PublishEnrollmentChange(c.Request.Context(), appCtx.GetLocalPubSub().GetUnblockPubSub(), updateCacheMsg); err != nil {
			logger.AppLogger.Error(
				c.Request.Context(),
				"cannot publish update user cache message",
				zap.Error(common.ErrInternal(err)),
				zap.Any("updateCacheMsg", updateCacheMsg),
			)
		}

		c.Next()
	}
}

func getVideoCacheInfo(
	c context.Context,
	db *gorm.DB,
	appCache component.AppCache,
	courseSlug string,
	videoId string,
) (*cache.VideoCacheInfo, error) {
	videoCacheInfo, err := appCache.GetVideoCache(c, courseSlug, videoId)
	if err != nil {
		return nil, err
	}

	if videoCacheInfo == nil {
		uid, err := common.FromBase58(videoId)
		if err != nil {
			return nil, err
		}
		videoId := uid.GetLocalID()

		videoStore := videostore.NewSQLStore(db)
		video, err := videoStore.FindOne(c, map[string]interface{}{"id": videoId})
		if err != nil {
			return nil, err
		}

		if err := appCache.SetVideoCache(c, courseSlug, *video); err != nil {
			return nil, err
		}
	}

	videoCacheInfo, err = appCache.GetVideoCache(c, courseSlug, videoId)
	if err != nil {
		return nil, err
	}

	return videoCacheInfo, nil
}
