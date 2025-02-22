package logger

import (
	"context"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otel_log "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func NewLogger() *logrus.Logger {
	logExporter, err := otlploghttp.New(context.Background(),
		otlploghttp.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlploghttp.WithInsecure(),
	)

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// create log provider
	log_provider := otel_log.NewLoggerProvider(
		otel_log.WithProcessor(
			otel_log.NewBatchProcessor(logExporter),
		),
	)

	defer log_provider.Shutdown(context.Background())

	// Create an *otellogrus.Hook and use it in your application.
	hook := otellogrus.NewHook("<signoz-golang>", otellogrus.WithLoggerProvider(log_provider))
	logger := logrus.New()
	logger.AddHook(hook)

	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "error" {
		logger.SetLevel(logrus.ErrorLevel)
	} else if logLevel == "info" {
		logger.SetLevel(logrus.InfoLevel)
	} else if logLevel == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set the newly created hook as a global logrus hook

	return logger
}

func GetTracerAndSpanID() (oteltrace.Tracer, oteltrace.Span) {
	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		),
	)

	tracer := otel.GetTracerProvider().Tracer("signoz-tracer")
	_, span := tracer.Start(context.Background(), "signoz-tracer")
	defer span.End()

	return tracer, span
}

func LogrusFields(span oteltrace.Span) logrus.Fields {
	spanCtx := span.SpanContext()
	fields := logrus.Fields{
		"trace_id": spanCtx.TraceID().String(),
		"span_id":  spanCtx.SpanID().String(),
	}
	return fields
}
