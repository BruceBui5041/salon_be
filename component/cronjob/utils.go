package cronjob

import (
	"context"
	"fmt"
	"salon_be/component/logger"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func MinutesToCronFormat(minutes int) (string, error) {
	if minutes < 1 {
		return "", fmt.Errorf("interval must be at least 1 minute")
	}

	if minutes < 60 {
		// For intervals less than an hour, use the format "@every Xm"
		return fmt.Sprintf("@every %dm", minutes), nil
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	if remainingMinutes == 0 {
		// For whole hour intervals, use the format "@every Xh"
		return fmt.Sprintf("@every %dh", hours), nil
	}

	// For intervals with both hours and minutes, construct the cron expression
	return fmt.Sprintf("*/%s */%s * * *",
		strconv.Itoa(remainingMinutes),
		strconv.Itoa(hours)), nil
}

func WithDurationLogging(ctx context.Context, jobName string, job func()) func() {
	return func() {
		tracer := otel.Tracer("CRONJOB duration track")
		ctx, span := tracer.Start(ctx, "duration track", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		start := time.Now()
		job()
		duration := time.Since(start)
		logger.AppLogger.Info(ctx, fmt.Sprintf("Job '%s' took %v to execute", jobName, duration))
	}
}
