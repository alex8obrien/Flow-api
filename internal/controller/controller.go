package controller

import (
	"github.com/alex8obrien/echo-v5-template/internal/telemetry"

	"github.com/labstack/echo/v5"
)

type UserController struct {
	uc     UserService
	tracer telemetry.Tracer
}

func NewController(uc UserService, tracer telemetry.Tracer) *UserController {
	return &UserController{uc: uc, tracer: tracer}
}

func (h *UserController) Routes(e *echo.Echo) {
	e.GET("/health", h.Health)
	e.POST("/users", h.Create)
	e.GET("/debug", h.Debug)
}
