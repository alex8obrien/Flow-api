package logger_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/alex8obrien/echo-v5-template/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Debug(t *testing.T) {
	// Define buffer
	var buf bytes.Buffer
	logger.New(slog.LevelInfo, &buf)

	// Test 1
	slog.DebugContext(t.Context(), "hello")
	output := strings.TrimSpace(buf.String())
	assert.Empty(t, output)
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger.New(slog.LevelInfo, &buf)

	slog.InfoContext(t.Context(), "hello")
	output := strings.TrimSpace(buf.String())
	assert.NotEmpty(t, output)
}
