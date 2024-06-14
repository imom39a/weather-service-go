package api

import (
	"net/http"

	"go.uber.org/zap"
)

/*
The composite handler is used to combine multiple handlers into a single handler.
For example, if you have a weather handler and a health handler, you can combine them into a single composite handler.
*/

type CompositeHandler struct {
	*WeatherHandler
	*HealthHandler
}

func NewCompositeHandler(logger *zap.Logger, httpClient *http.Client) *CompositeHandler {
	return &CompositeHandler{
		WeatherHandler: NewWeatherHandler(logger, httpClient),
		HealthHandler:  NewHealthHandler(logger),
	}
}
