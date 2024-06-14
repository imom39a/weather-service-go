package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	errors "weathe-service/common/error"
	"weathe-service/internal/api"
	"weathe-service/internal/entity"

	"github.com/eapache/go-resiliency/breaker"
	"go.uber.org/zap"
)

type WeatherService struct {
	logger  *zap.Logger
	apiKey  string
	breaker *breaker.Breaker
}

func NewWeatherService(logger *zap.Logger) *WeatherService {
	key := os.Getenv("WEATHER_SERVICE_API_KEY")
	if key == "" {
		log.Fatal("WEATHER_SERVICE_API_KEY is not set")
		panic("WEATHER_SERVICE_API_KEY is not set")
	}
	return &WeatherService{
		apiKey:  key,
		logger:  logger,
		breaker: breaker.New(3, 1, 5), // 3 failures in 5 seconds will trip the breaker
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float32, unit string) (*api.WeatherResponse, error) {
	s.logger.Debug("GetWeather service called with lat and lon", zap.Float32("lat", lat), zap.Float32("lon", lon))
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=%s",
		strconv.FormatFloat(float64(lat), 'f', 6, 32),
		strconv.FormatFloat(float64(lon), 'f', 6, 32),
		s.apiKey,
		unit,
	)
	var weatherResponse entity.WeatherResponse
	err := s.breaker.Run(func() error {
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			s.logger.Error("failed to create request", zap.Error(err))
			return errors.NewAPIError(http.StatusInternalServerError, "failed to create request", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				s.logger.Error("network error", zap.Error(err))
				return errors.NewAPIError(http.StatusServiceUnavailable, "network error, unable to reach weather service", err)
			}
			s.logger.Error("failed to fetch weather data", zap.Error(err))
			return errors.NewAPIError(http.StatusInternalServerError, "failed to fetch weather data", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s.logger.Error("response code is not 200", zap.Int("statusCode", resp.StatusCode))
			return errors.NewAPIError(resp.StatusCode, "failed to fetch weather data. response code is not 200", nil)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Error("failed to read weather data", zap.Error(err))
			return errors.NewAPIError(http.StatusInternalServerError, "Error processing weather data", err)
		}

		err = json.Unmarshal(body, &weatherResponse)
		if err != nil {
			s.logger.Error("failed to unmarshal weather data", zap.Error(err))
			return errors.NewAPIError(http.StatusInternalServerError, "Error processing weather data", err)
		}

		return nil
	})

	if err != nil {
		if err == breaker.ErrBreakerOpen {
			s.logger.Warn("circuit breaker triggered")
			return nil, errors.NewAPIError(http.StatusServiceUnavailable, "service unavailable, circuit breaker triggered", err)
		} else {
			s.logger.Error("Error processing weather data", zap.Error(err))
			return nil, err
		}
	}

	tempDescription := getTemperatureDescription(weatherResponse.Main.Temp)

	weatherCondition := "unknown"
	if len(weatherResponse.Weather) > 0 {
		weatherCondition = weatherResponse.Weather[0].Main
	}

	return &api.WeatherResponse{
		TemperatureCondition: &tempDescription,
		WeatherCondition:     &weatherCondition,
		Temperature:          float64ToFloat32Ptr(weatherResponse.Main.Temp),
		Unit:                 &unit,
	}, nil
}

func float64ToFloat32Ptr(f float64) *float32 {
	if f == 0 {
		return nil
	}
	f32 := float32(f)
	return &f32
}

func getTemperatureDescription(temp float64) string {
	switch {
	case temp < 283.15:
		return "cold"
	case temp < 293.15:
		return "moderate"
	case temp < 303.15:
		return "warm"
	default:
		return "hot"
	}
}
