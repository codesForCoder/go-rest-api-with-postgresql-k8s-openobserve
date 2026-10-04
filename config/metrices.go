package config

import (
	"context"
	"log"
	"os"
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

	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(15*time.Second),
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
