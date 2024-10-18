package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var AppLogger *Logger

func CreateAppLogger(ctx context.Context) {
	logLevel := getLogLevelFromEnv()

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)

	// Modify the EncoderConfig to change how the caller is captured
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Create a new development config
	devConfig := zap.NewDevelopmentConfig()
	devConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Combine the production and development configs
	opts := []zap.Option{
		zap.AddCallerSkip(1), // Skip the wrapper logger calls
		zap.AddCaller(),      // Add caller information
	}

	logger, err := config.Build(opts...)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	AppLogger = &Logger{logger}

	AppLogger.Info(ctx, "Logger initialized", zap.String("level", logLevel.String()))
}

func getLogLevelFromEnv() zapcore.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	levelStr = strings.ToUpper(levelStr)

	switch levelStr {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "DPANIC":
		return zapcore.DPanicLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel // Default to Info if not specified or invalid
	}
}

func (l *Logger) UpdateLogLevel(ctx context.Context, level zapcore.Level) {
	l.Core().Enabled(level)
	l.Info(ctx, "Log level updated", zap.String("newLevel", level.String()))
}

func addOtelFields(ctx context.Context, fields []zap.Field) []zap.Field {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return append(fields, zap.String("trace_id", traceID), zap.String("span_id", spanID))
}

// Debug logs a message at DebugLevel with OpenTelemetry context
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, addOtelFields(ctx, fields)...)
}

// Info logs a message at InfoLevel with OpenTelemetry context
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, addOtelFields(ctx, fields)...)
}

// Warn logs a message at WarnLevel with OpenTelemetry context
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, addOtelFields(ctx, fields)...)
}

// Error logs a message at ErrorLevel with OpenTelemetry context
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Error(msg, addOtelFields(ctx, fields)...)
}

// DPanic logs a message at DPanicLevel with OpenTelemetry context
func (l *Logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.DPanic(msg, addOtelFields(ctx, fields)...)
}

// Panic logs a message at PanicLevel with OpenTelemetry context
func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, addOtelFields(ctx, fields)...)
}

// Fatal logs a message at FatalLevel with OpenTelemetry context
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, addOtelFields(ctx, fields)...)
}

// WithContext returns a logger with the given context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{l.Logger.With(addOtelFields(ctx, []zap.Field{})...)}
}
