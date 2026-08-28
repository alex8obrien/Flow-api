package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alex8obrien/flow/internal/controller"
	"github.com/alex8obrien/flow/internal/db"
	"github.com/alex8obrien/flow/internal/environment"
	"github.com/alex8obrien/flow/internal/logger"
	"github.com/alex8obrien/flow/internal/repository"
	"github.com/alex8obrien/flow/internal/server"
	"github.com/alex8obrien/flow/internal/service"
	"github.com/alex8obrien/flow/internal/telemetry/otel"
)

const timeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("start up failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Setup environment
	if err := environment.InitEnv(".env"); err != nil {
		return err
	}

	// Setup logger
	logger.New(slog.LevelDebug, os.Stdout)

	// Setup context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Setup Otel
	otelCfg, err := otel.InitialiseOtel(os.Getenv(environment.ApplicationName))
	if err != nil {
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if shutdownErr := otelCfg.Shutdown(shutdownCtx); shutdownErr != nil {
			slog.Error("otel shutdown failed", "error", shutdownErr)
		}
	}()

	// Define pool
	pool, err := db.NewPool(ctx, os.Getenv(environment.Dsn))
	if err != nil {
		return fmt.Errorf("failed to create db pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	defer pool.Close()

	// Setup application controler and server configuration
	linkRepoTracer := otelCfg.NewTracer("link-repository")
	linkServiceTracer := otelCfg.NewTracer("link-service")
	linkControllerTracer := otelCfg.NewTracer("link-controller")
	r := repository.NewLinkRepository(pool, linkRepoTracer)
	s := service.NewLinkService(r, linkServiceTracer)
	c := controller.NewController(s, linkControllerTracer)
	srv := server.Config{
		OtelMiddleware: otel.Middleware(otelCfg),
		Address:        os.Getenv(environment.Address),
		Origins:        []string{os.Getenv(environment.AllowOrigins)},
		CertPath:       os.Getenv(environment.CertFile),
		KeyPath:        os.Getenv(environment.KeyFile),
	}

	// Create Echo server, define routes and start
	e := server.InitServer(srv)
	c.Routes(e)
	return srv.Start(ctx, e)
}
