package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alex8obrien/flow/internal/controller"
	"github.com/alex8obrien/flow/internal/server"
	"github.com/alex8obrien/flow/internal/test/telemetry"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	e = echo.New()
	e.Validator = server.NewCustomValidator()
	controller := controller.NewController(NewMockService(), telemetry.NewTracer(), testAuthMiddleware)
	controller.Routes(e)
	os.Exit(m.Run())
}

func testAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) != "Bearer test-token" {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return next(c)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body)
}

func TestCreateEndpoint(t *testing.T) {
	jsonBody := []byte(
		`{"destination_url": "https://example.com", "display_name": "Example", "description": "My description"}`,
	)
	body := bytes.NewReader(jsonBody)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, rec.Body)
}

func TestCreateEndpoint_NullDescription(t *testing.T) {
	jsonBody := []byte(
		`{"destination_url": "https://example.com", "display_name": "Example"}`,
	)
	body := bytes.NewReader(jsonBody)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, rec.Body)
}

func TestRegisterEndpoint_UnprocessableDestinationURL(t *testing.T) {
	jsonBody := []byte(`{"destination_url": "https://", "display_name": "Example"}`)
	body := bytes.NewReader(jsonBody)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.NotEmpty(t, rec.Body)
}

func TestCreateEndpoint_Unauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewReader([]byte(`{}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
