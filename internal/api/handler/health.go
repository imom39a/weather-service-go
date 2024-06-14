package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) GetHealth(ctx echo.Context) error {
	h.logger.Debug("Health check handler")
	return ctx.String(http.StatusOK, "Up and running!")
}
