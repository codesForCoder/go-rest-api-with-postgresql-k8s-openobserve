package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func SetupMetrics() func() {
	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpoint(
			os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		),
		otlpmetrichttp.WithInsecure(),
	)

	if err != nil {
		panic(err)
	}
	// Default: 15 seconds
	interval := 15000
	if value := os.Getenv("OTEL_METRIC_EXPORT_INTERVAL"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			log.Printf("Invalid OTEL_METRIC_EXPORT_INTERVAL=%q; using default 15000ms", value)
		} else {
			interval = parsed
		}
	}

	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(time.Duration(interval)*time.Millisecond),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(mp)

	return func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := mp.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}
}
