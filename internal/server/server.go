package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	bodyLimit           int64         = 2 * middleware.MB // 2 MB
	hstsMaximumAge      int           = 31536000          // 1 Year
	gracefulTimeout     time.Duration = 10 * time.Second  // 10 Seconds
	accessControlMaxAge int           = 3600              // 1 Hour
	rateLimit           float64       = 20                // 20 Requests / Second
)

type Config struct {
	OtelMiddleware echo.MiddlewareFunc
	Address        string
	Origins        []string
	KeyPath        string
	CertPath       string
}

func InitServer(cfg Config) *echo.Echo {
	e := echo.New()
	e.Validator = NewCustomValidator()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger()) // TODO: Replace with RequestLoggerWithConfig
	e.Use(cfg.OtelMiddleware)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            hstsMaximumAge,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	// TODO: Add rate limit middleware
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rateLimit)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Origins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodOptions, http.MethodGet, http.MethodPost},
		MaxAge:       accessControlMaxAge,
	}))
	// TODO: Add CSRF
	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{}))
	e.Use(middleware.BodyLimit(bodyLimit)) // 2MB body limit
	e.Use(middleware.Gzip())

	return e
}

func (cfg Config) Start(ctx context.Context, e *echo.Echo) error {
	cert, err := os.ReadFile(cfg.CertPath)
	if err != nil {
		e.Logger.ErrorContext(ctx, "failed to read cert file", "error", err)
		return err
	}

	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		e.Logger.ErrorContext(ctx, "failed to read key file", "error", err)
		return err
	}

	sc := echo.StartConfig{
		Address:         cfg.Address,
		GracefulTimeout: gracefulTimeout,
	}

	// Start server
	if err = sc.StartTLS(ctx, e, cert, key); err != nil {
		e.Logger.ErrorContext(ctx, "server stopped", "error", err)
		return err
	}

	return nil
}
