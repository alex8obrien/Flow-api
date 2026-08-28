package telemetry

import "context"

type Span interface {
	End()
	AddEvent(name string, attrs ...Attribute)
	RecordError(err error)
	SetAttributes(attrs ...Attribute)
	SetStatus(status Status, description string)
}

type Tracer interface {
	Start(ctx context.Context, name string, attrs ...Attribute) (context.Context, Span)
}
