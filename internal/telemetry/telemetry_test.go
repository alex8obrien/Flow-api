package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, Attribute{Key: "http.method", Value: "GET"}, String("http.method", "GET"))
}
