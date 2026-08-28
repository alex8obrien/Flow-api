package telemetry

import (
	"context"

	"github.com/alex8obrien/echo-v5-template/internal/telemetry"
)

type Tracer struct{}

type Span struct{}

func NewTracer() telemetry.Tracer {
	return Tracer{}
}

func (Tracer) Start(ctx context.Context, _ string, _ ...telemetry.Attribute) (context.Context, telemetry.Span) {
	return ctx, Span{}
}

func (Span) End() {}

func (Span) AddEvent(_ string, _ ...telemetry.Attribute) {}

func (Span) RecordError(_ error) {}

func (Span) SetAttributes(_ ...telemetry.Attribute) {}

func (Span) SetStatus(_ telemetry.Status, _ string) {}
