package service

import "context"

type WeatherService struct {
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float32) string {
	// call weather api here
	return "Sunny"
}
