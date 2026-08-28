package otel

import (
	"github.com/alex8obrien/flow/internal/telemetry"

	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func InitialiseOtel(applicationName string) (*TraceConfig, error) {
	// Initialise tracer provider
	tp, err := initTracer()
	if err != nil {
		return nil, err
	}

	return &TraceConfig{
		ApplicationName: applicationName,
		Tp:              tp,
	}, nil
}

func Middleware(cfg *TraceConfig) echo.MiddlewareFunc {
	// Return middleware or error
	return echootel.NewMiddlewareWithConfig(echootel.Config{
		ServerName:     cfg.ApplicationName,
		TracerProvider: cfg.Tp,
	})
}

func toOTelAttributes(attrs []telemetry.Attribute) []attribute.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	otelAttrs := make([]attribute.KeyValue, 0, len(attrs))
	for _, a := range attrs {
		otelAttrs = append(otelAttrs, attribute.String(a.Key, a.Value))
	}

	return otelAttrs
}

func toOTelStatusCode(status telemetry.Status) codes.Code {
	switch status {
	case telemetry.StatusSuccess:
		return codes.Ok
	case telemetry.StatusError:
		return codes.Error
	case telemetry.StatusUnset:
		return codes.Unset
	default:
		return codes.Unset
	}
}
