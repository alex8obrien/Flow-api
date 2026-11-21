package controller

import (
	"github.com/alex8obrien/flow/internal/telemetry"

	"github.com/labstack/echo/v5"
)

type LinkController struct {
	srv      LinkService
	tracer   telemetry.Tracer
	authware echo.MiddlewareFunc
}

func NewController(srv LinkService, tracer telemetry.Tracer, authware echo.MiddlewareFunc) *LinkController {
	return &LinkController{srv: srv, tracer: tracer, authware: authware}
}

func (h *LinkController) Routes(e *echo.Echo) {
	e.GET("/livez", h.Livez)
	e.GET("/readyz", h.Readyz)
	e.POST("/links", h.Create, h.authware)
	e.GET("/:id", h.Redirect)
}
