package api

import (
	"net/http"
	"weathe-service/internal/api"
	"weathe-service/internal/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WeatherHandler struct {
	logger         *zap.Logger
	weatherService *service.WeatherService
}

func NewWeatherHandler(logger *zap.Logger, httpClient *http.Client) *WeatherHandler {
	return &WeatherHandler{
		logger:         logger,
		weatherService: service.NewWeatherService(logger, httpClient),
	}
}
func (h *WeatherHandler) GetWeather(ctx echo.Context, params api.GetWeatherParams) error {
	h.logger.Debug("GetWeather handler")

	unit := "imperial"
	if params.Unit != nil {
		unit = string(*params.Unit)
	}

	res, err := h.weatherService.GetWeather(ctx.Request().Context(), params.Lat, params.Lon, unit)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}
