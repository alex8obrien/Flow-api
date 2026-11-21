package controller

import (
	"net/http"

	"github.com/alex8obrien/flow/internal/domain"
	"github.com/alex8obrien/flow/internal/telemetry"

	"github.com/labstack/echo/v5"
)

func (ctrl *LinkController) Livez(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (ctrl *LinkController) Readyz(c *echo.Context) error {
	ctx, span := ctrl.tracer.Start(c.Request().Context(), "controller.links.readyz")
	defer span.End()
	span.AddEvent("controller.links.readyz.request_received")

	r := ctrl.srv.Health(ctx)

	// Define response
	code := http.StatusOK
	if r.Overall() != domain.OK {
		code = http.StatusServiceUnavailable
	}

	// Return response
	return c.JSON(code, r)
}

func (ctrl *LinkController) Create(c *echo.Context) error {
	ctx, span := ctrl.tracer.Start(c.Request().Context(), "controller.links.create")
	defer span.End()
	span.AddEvent("controller.links.create.request_received")

	var dto CreateLinkDTO
	if err := c.Bind(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "bind failed")
		span.AddEvent("controller.links.create.bind_failed")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	span.AddEvent("controller.links.create.bind_succeeded")

	if err := c.Validate(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "validate failed")
		span.AddEvent("controller.links.create.validate_failed")
		return err
	}

	span.AddEvent("controller.links.create.validate_succeeded")
	link, err := ctrl.srv.Create(ctx, dto.DestinationURL, dto.DisplayName, dto.Description)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "create failed")
		span.AddEvent("controller.links.create.create_failed")
		return err
	}

	span.SetAttributes(telemetry.String("link.id", link.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "link created")
	span.AddEvent("controller.links.create.completed")
	return c.JSON(http.StatusCreated, link)
}

func (ctrl *LinkController) Redirect(c *echo.Context) error {
	ctx, span := ctrl.tracer.Start(c.Request().Context(), "controller.links.redirect")
	defer span.End()
	span.AddEvent("controller.links.redirect.request_received")

	id := c.Param("id")
	link, err := ctrl.srv.Redirect(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, "redirect failed")
		span.AddEvent("controller.links.redirect.redirect_failed")
		return err
	}

	span.SetStatus(telemetry.StatusSuccess, "redirect started")
	span.AddEvent("controller.links.redirect.completed")
	return c.Redirect(http.StatusMovedPermanently, link.String())
}
