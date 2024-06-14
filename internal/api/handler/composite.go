package api

type CompositeHandler struct {
	*WeatherHandler
	*HealthHandler
}

func NewCompositeHandler() *CompositeHandler {
	return &CompositeHandler{
		WeatherHandler: NewWeatherHandler(),
		HealthHandler:  NewHealthHandler(),
	}
}
