package controller

import (
	"github.com/alex8obrien/echo-v5-template/internal/telemetry"

	"github.com/labstack/echo/v5"
)

type LinkController struct {
	srv    LinkService
	tracer telemetry.Tracer
}

func NewController(srv LinkService, tracer telemetry.Tracer) *LinkController {
	return &LinkController{srv: srv, tracer: tracer}
}

func (h *LinkController) Routes(e *echo.Echo) {
	e.GET("/health", h.Health)
	e.POST("/links", h.Create)
	e.GET("/:id", h.Redirect)
}
