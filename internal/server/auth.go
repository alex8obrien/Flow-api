package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/alex8obrien/flow/internal/auth"

	"github.com/labstack/echo/v5"
)

const authClaimsContextKey = "auth.claims"

type TokenVerifier interface {
	Verify(context.Context, string) (*auth.Claims, error)
}

func AuthMiddleware(verifier TokenVerifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get(echo.HeaderAuthorization)
			if !strings.HasPrefix(header, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}

			claims, err := verifier.Verify(
				c.Request().Context(),
				strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")),
			)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
			}
			c.Set(authClaimsContextKey, claims)
			return next(c)
		}
	}
}
