package otel

import (
	"context"

	"github.com/alex8obrien/flow/internal/telemetry"

	echootel "github.com/labstack/echo-opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type TraceConfig struct {
	ApplicationName string
	Tp              *sdktrace.TracerProvider
}

const TracerKey = echootel.TracerKey

func (cfg *TraceConfig) Tracer(name string) trace.Tracer {
	return cfg.Tp.Tracer(name)
}

func (cfg *TraceConfig) NewTracer(name string) telemetry.Tracer {
	return otelTracer{tracer: cfg.Tp.Tracer(name)}
}

func (cfg *TraceConfig) Shutdown(ctx context.Context) error {
	return cfg.Tp.Shutdown(ctx)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}

type otelTracer struct {
	tracer trace.Tracer
}

func (t otelTracer) Start(
	ctx context.Context,
	name string,
	attrs ...telemetry.Attribute,
) (context.Context, telemetry.Span) {
	ctx, span := t.tracer.Start(ctx, name, trace.WithAttributes(toOTelAttributes(attrs)...))
	return ctx, otelSpan{span: span}
}

type otelSpan struct {
	span trace.Span
}

func (s otelSpan) End() {
	s.span.End()
}

func (s otelSpan) AddEvent(name string, attrs ...telemetry.Attribute) {
	s.span.AddEvent(name, trace.WithAttributes(toOTelAttributes(attrs)...))
}

func (s otelSpan) RecordError(err error) {
	s.span.RecordError(err)
}

func (s otelSpan) SetAttributes(attrs ...telemetry.Attribute) {
	s.span.SetAttributes(toOTelAttributes(attrs)...)
}

func (s otelSpan) SetStatus(status telemetry.Status, description string) {
	s.span.SetStatus(toOTelStatusCode(status), description)
}
