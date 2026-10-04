package config

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type OTelHandler struct {
	slog.Handler
}

func (h *OTelHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {

	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	if spanContext.IsValid() {
		record.AddAttrs(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, record)
}

func SetupLogger() {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)

	otelHandler := &OTelHandler{
		Handler: handler,
	}

	slog.SetDefault(slog.New(otelHandler))
}
