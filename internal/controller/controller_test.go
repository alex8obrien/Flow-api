package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alex8obrien/echo-v5-template/internal/controller"
	"github.com/alex8obrien/echo-v5-template/internal/server"
	"github.com/alex8obrien/echo-v5-template/internal/test/telemetry"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// TODO: Design better echo integration with controller test

var e *echo.Echo

func TestMain(m *testing.M) {
	e = echo.New()
	e.Validator = server.NewCustomValidator()
	controller := controller.NewController(NewMockService(), telemetry.NewTracer())
	controller.Routes(e)
	os.Exit(m.Run())
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body)
}

func TestRequestEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/debug", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body)
}

func TestRegisterEndpoint(t *testing.T) {
	jsonBody := []byte(`{"name": "testuser", "email": "test@example.com"}`)
	body := bytes.NewReader(jsonBody)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, rec.Body)
}

func TestRegisterEndpoint_UnprocessableEmail(t *testing.T) {
	jsonBody := []byte(`{"name": "testuser", "email": "broken.com"}`)
	body := bytes.NewReader(jsonBody)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.NotEmpty(t, rec.Body)
}
