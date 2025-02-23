package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otel_log "go.opentelemetry.io/otel/sdk/log"
)

// NewLogger creates and returns a configured Logrus logger with OpenTelemetry hook
func NewLogger() (*logrus.Logger, error) {
	logger := logrus.New()

	return logger, nil
}

func GetLoggerWithContext(ctx context.Context) (*logrus.Logger, error, *otel_log.LoggerProvider) {
	logger, err := NewLogger()
	if err != nil {
		return nil, err, nil
	}

	exporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("failed to create OpenTelemetry exporter: %v", err), nil
	}

	// Create an OpenTelemetry log processor

	logProvider := otel_log.NewLoggerProvider(
		otel_log.WithProcessor(
			otel_log.NewBatchProcessor(exporter),
		),
	)

	// Set the OpenTelemetry log processor

	hook := otellogrus.NewHook("<signoz-golang>", otellogrus.WithLoggerProvider(logProvider), otellogrus.WithLevels(logrus.AllLevels))

	logger.AddHook(hook)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger, nil, logProvider
}

func AddHooks(logger *logrus.Logger, hooks ...logrus.Hook) {
	for _, hook := range hooks {
		logger.AddHook(hook)
	}
}
