package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alex8obrien/flow/internal/controller"
	"github.com/alex8obrien/flow/internal/domain"
	"github.com/alex8obrien/flow/internal/server"
	"github.com/alex8obrien/flow/internal/test/telemetry"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func newReadyzTestEcho(readyz *domain.Readyz) *echo.Echo {
	testEcho := echo.New()
	controller.NewController(
		MockService{readyz: readyz},
		telemetry.NewTracer(),
		testAuthMiddleware,
	).Routes(testEcho)
	return testEcho
}

func TestLivezEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		httpMethod         string
		expectedStatusCode int
	}{
		{name: "GET is alive", httpMethod: http.MethodGet, expectedStatusCode: http.StatusOK},
		{name: "POST is not allowed", httpMethod: http.MethodPost, expectedStatusCode: http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.httpMethod, "/livez", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatusCode, rec.Code)
			if test.httpMethod == http.MethodGet {
				assert.Empty(t, rec.Body.String())
				assert.Empty(t, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

func TestReadyzEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		readyz             *domain.Readyz
		expectedStatusCode int
		expectedOverall    string
		expectedDB         string
	}{
		{
			name: "healthy",
			readyz: func() *domain.Readyz {
				r := domain.NewReadyz()
				r.DBChecked()
				return r
			}(),
			expectedStatusCode: http.StatusOK,
			expectedOverall:    "ok",
			expectedDB:         "ok",
		},
		{
			name:               "database unavailable",
			readyz:             domain.NewReadyz(),
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedOverall:    "down",
			expectedDB:         "down",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			rec := httptest.NewRecorder()
			newReadyzTestEcho(test.readyz).ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatusCode, rec.Code)
			assert.Equal(t, echo.MIMEApplicationJSON, rec.Header().Get(echo.HeaderContentType))

			var body map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, map[string]any{
				"version": domain.Version,
				"overall": test.expectedOverall,
				"db":      test.expectedDB,
			}, body)
		})
	}
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

func TestCreateEndpoint_UnprocessableDestinationURL(t *testing.T) {
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
