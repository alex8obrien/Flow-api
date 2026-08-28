package controller

import (
	"fmt"
	"net/http"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
	"github.com/alex8obrien/echo-v5-template/internal/telemetry"

	"github.com/labstack/echo/v5"
)

func (h *UserController) Health(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *UserController) Create(c *echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "controller.user.create")
	defer span.End()
	span.AddEvent("controller.user.create.request_received")

	var dto CreateUserDTO
	if err := c.Bind(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "bind failed")
		span.AddEvent("controller.user.create.bind_failed")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	span.AddEvent("controller.user.create.bind_succeeded")
	span.SetAttributes(telemetry.String("user.email_domain", dto.Email))

	if err := c.Validate(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "validate failed")
		span.AddEvent("controller.user.create.validate_failed")
		return err
	}

	span.AddEvent("controller.user.create.validate_succeeded")
	user, err := h.uc.Register(ctx, dto.Name, dto.Email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "register failed")
		span.AddEvent("controller.user.create.register_failed")
		return err
	}

	span.SetAttributes(telemetry.String("user.id", user.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "user created")
	span.AddEvent("controller.user.create.completed")
	return c.JSON(http.StatusCreated, user)
}

func (h *UserController) Debug(c *echo.Context) error {
	req := c.Request()

	tlsInfo := "NO TLS"

	if req.TLS != nil {
		tlsInfo = fmt.Sprintf(
			"TLS=%d, ALPN=%q",
			req.TLS.Version,
			req.TLS.NegotiatedProtocol,
		)
	}

	format := `
    <code>
      Version: %s<br>
      Protocol: %s<br>
      TLS: %s<br>
      Host: %s<br>
      Remote Address: %s<br>
      Method: %s<br>
      Path: %s<br>
    </code>
    `

	return c.HTML(http.StatusOK, fmt.Sprintf(
		format,
		domain.Version,
		req.Proto,
		tlsInfo,
		req.Host,
		req.RemoteAddr,
		req.Method,
		req.URL.Path,
	))
}
