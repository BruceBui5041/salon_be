package job

import (
	"context"
	"time"
	"video_server/component"
	"video_server/component/logger"
	models "video_server/model"
	"video_server/model/course/coursestore"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func CourseUpdateCountFieldJob(ctx context.Context, appCtx component.AppContext) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		tracer := otel.Tracer("CRONJOB")
		ctx, span := tracer.Start(ctx, "cron job update course count field", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)

		courses, err := courseStore.FindAll(
			ctx,
			map[string]interface{}{},
			"Lessons",
			"Enrollments.Payment",
			"Comments",
		)

		if err != nil {
			logger.AppLogger.Error(ctx, "failed to get courses", zap.Error(err))
			return
		}

		for _, course := range courses {
			var studentCount uint16
			for _, enroll := range course.Enrollments {
				if enroll.Payment.TransactionStatus == "completed" {
					studentCount += 1
				}
			}

			reviewInfo := calculateReviewInfo(course.Comments)
			avgRating := calculateAverageRating(reviewInfo)
			ratingCount := calculateTotalRatingCount(reviewInfo)

			updateCourse := models.Course{
				LessonCount:   uint16(len(course.Lessons)),
				StudentCount:  studentCount,
				ReviewInfo:    reviewInfo,
				AverageRating: avgRating,
				RatingCount:   ratingCount,
			}

			if err := courseStore.Update(ctx, course.Id, &updateCourse); err != nil {
				logger.AppLogger.Error(ctx, "failed to update courses", zap.Error(err))
				return
			}
		}
	}
}

func calculateReviewInfo(comments []*models.Comment) models.ReviewInfos {
	reviewInfo := make(models.ReviewInfos, 5)
	for i := range reviewInfo {
		reviewInfo[i] = models.ReviewInfo{Stars: uint8(i + 1), Count: 0}
	}

	for _, comment := range comments {
		if comment.Rate >= 1 && comment.Rate <= 5 {
			reviewInfo[comment.Rate-1].Count++
		}
	}

	return reviewInfo
}

func calculateAverageRating(reviewInfo models.ReviewInfos) decimal.Decimal {
	totalRating := decimal.Zero
	totalCount := uint(0)

	for _, info := range reviewInfo {
		rating := decimal.NewFromInt(int64(info.Stars))
		count := decimal.NewFromInt(int64(info.Count))
		totalRating = totalRating.Add(rating.Mul(count))
		totalCount += info.Count
	}

	if totalCount == 0 {
		return decimal.Zero
	}

	return totalRating.Div(decimal.NewFromInt32(int32(totalCount))).Round(1)
}

func calculateTotalRatingCount(reviewInfo models.ReviewInfos) uint {
	var totalCount uint
	for _, info := range reviewInfo {
		totalCount += uint(info.Count)
	}
	return totalCount
}
