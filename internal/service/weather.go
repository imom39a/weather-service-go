package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"weathe-service/internal/entity"
)

type WeatherService struct {
	apiKey string
}

func NewWeatherService() *WeatherService {
	key := os.Getenv("WEATHER_SERVICE_API_KEY")
	if key == "" {
		log.Fatal("WEATHER_SERVICE_API_KEY is not set")
		panic("WEATHER_SERVICE_API_KEY is not set")
	}
	return &WeatherService{
		apiKey: key,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float32) string {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s",
		strconv.FormatFloat(float64(lat), 'f', 6, 32),
		strconv.FormatFloat(float64(lon), 'f', 6, 32),
		s.apiKey,
	)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to fetch weather data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weatherResponse entity.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatal(err)
	}

	return weatherResponse.Weather[0].Main
}
