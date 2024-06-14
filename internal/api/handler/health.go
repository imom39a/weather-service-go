package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Up and running!")
}
