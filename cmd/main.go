package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otel_log "go.opentelemetry.io/otel/sdk/log"
)

func main() {
	logExporter, err := otlploghttp.New(context.Background())
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
	hook := otellogrus.NewHook("<signoz-golang>", otellogrus.WithLoggerProvider(log_provider), otellogrus.WithLevels(logrus.AllLevels))

	// Set the newly created hook as a global logrus hook
	logrus.AddHook(hook)

	// add a debug log

	// logrus.Warn("This is a warning log")
	logrus.Trace("This is a trace log")
}
