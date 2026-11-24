package telemetry_test

import (
	"testing"

	"github.com/alex8obrien/flow/internal/telemetry"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, telemetry.Attribute{Key: "http.method", Value: "GET"}, telemetry.String("http.method", "GET"))
}
