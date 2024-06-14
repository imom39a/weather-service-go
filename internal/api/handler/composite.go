package api

import "go.uber.org/zap"

type CompositeHandler struct {
	*WeatherHandler
	*HealthHandler
}

func NewCompositeHandler(logger *zap.Logger) *CompositeHandler {
	return &CompositeHandler{
		WeatherHandler: NewWeatherHandler(logger),
		HealthHandler:  NewHealthHandler(logger),
	}
}
