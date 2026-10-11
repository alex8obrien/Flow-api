package server_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex8obrien/flow/internal/auth"
	"github.com/alex8obrien/flow/internal/server"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type testTokenVerifier struct {
	claims *auth.Claims
	err    error
}

func (v testTokenVerifier) Verify(context.Context, string) (*auth.Claims, error) {
	return v.claims, v.err
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		verifier   server.TokenVerifier
		statusCode int
	}{
		{name: "missing token", statusCode: http.StatusUnauthorized},
		{
			name:       "invalid token",
			header:     "Bearer invalid",
			verifier:   testTokenVerifier{err: errors.New("invalid")},
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "valid token",
			header:     "Bearer valid",
			verifier:   testTokenVerifier{claims: &auth.Claims{}},
			statusCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/protected", func(c *echo.Context) error {
				return c.NoContent(http.StatusNoContent)
			}, server.AuthMiddleware(test.verifier))

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(echo.HeaderAuthorization, test.header)
			res := httptest.NewRecorder()
			e.ServeHTTP(res, req)

			assert.Equal(t, test.statusCode, res.Code)
		})
	}
}
