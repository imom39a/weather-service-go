package api

import (
	"net/http"
	"weathe-service/internal/api"
	"weathe-service/internal/service"

	"github.com/labstack/echo/v4"
)

type WeatherHandler struct {
	weatherService *service.WeatherService
}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{
		weatherService: service.NewWeatherService(),
	}
}
func (h *WeatherHandler) GetWeather(ctx echo.Context, params api.GetWeatherParams) error {
	res := h.weatherService.GetWeather(ctx.Request().Context(), params.Lat, params.Lon)

	return ctx.JSON(http.StatusOK, res)
}
