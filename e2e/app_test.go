package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"weathe-service/common/logger"
	"weathe-service/common/server"
	"weathe-service/internal/api"
	handlers "weathe-service/internal/api/handler"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	tests := []ApiTest{
		{
			name:        "Health check",
			description: "Should return status 200",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/health", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusOK, rec.StatusCode)
				assert.Equal(t, "Up and running!", string(body))
			},
		},
		{
			name:        "Get weather",
			description: "Should return status 200",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lat=1.0&lon=1.0&unit=metric", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusOK, rec.StatusCode)
				assert.NotEmpty(t, body)
			},
		},
		{
			name:        "Get weather with invalid lat and lon",
			description: "Should return status 400",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lat=invalid&lon=invalid&unit=metric", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusBadRequest, rec.StatusCode)
				assert.NotEmpty(t, body)
			},
		},
		{
			name:        "Get weather with invalid unit",
			description: "Should return status 400",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lat=1.0&lon=1.0&unit=invalid", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusBadRequest, rec.StatusCode)
				assert.NotEmpty(t, body)
			},
		},
		{
			name:        "Get weather with missing lat",
			description: "Should return status 400",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lon=1.0&unit=metric", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusBadRequest, rec.StatusCode)
				assert.NotEmpty(t, body)
			},
		},
		{
			name:        "Get weather with missing lon",
			description: "Should return status 400",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lat=1.0&unit=metric", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				assert.Equal(t, http.StatusBadRequest, rec.StatusCode)
				assert.NotEmpty(t, body)
			},
		},
		{
			name:        "Get weather with missing unit",
			description: "Should return with default unit as imperial",
			prepare: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/weather?lat=1.0&lon=1.0", nil)
				return req
			},
			assertFunc: func(t *testing.T, rec *http.Response, body []byte) {
				var res api.WeatherResponse
				_ = json.Unmarshal(body, &res)
				assert.Equal(t, http.StatusOK, rec.StatusCode)
				assert.NotEmpty(t, body)
				assert.Equal(t, "imperial", *res.Unit)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.prepare()
			rec := httptest.NewRecorder()
			e := server.CreateServer()
			e.Use(server.GetSwaggerValidatorMiddleware("../spec/weather-service.yaml"))

			api.RegisterHandlers(e, handlers.NewCompositeHandler(logger.NewLogger(), &http.Client{}))
			e.ServeHTTP(rec, req)
			body, _ := io.ReadAll(rec.Body)
			tt.assertFunc(t, rec.Result(), body)
		})
	}
}

func TestApiForCircuteBreakerScenarios(t *testing.T) {
	// Mock the HTTP client
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Simulate network failures to trigger the circuit breaker
	httpmock.RegisterResponder(http.MethodGet, "https://api.openweathermap.org/data/2.5/weather?lat=1.0&lon=1.0&appid=test",
		httpmock.NewErrorResponder(assert.AnError))

	httpClient := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "/weather?lat=1.0&lon=1.0&unit=metric", nil)
	rec := httptest.NewRecorder()
	e := server.CreateServer()
	e.Use(server.GetSwaggerValidatorMiddleware("../spec/weather-service.yaml"))
	api.RegisterHandlers(e, handlers.NewCompositeHandler(logger.NewLogger(), httpClient))
	e.ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Body)

	// Assertions
	assert.Equal(t, http.StatusServiceUnavailable, rec.Result().StatusCode)
	assert.NotEmpty(t, body)

	var res api.ErrorResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert the circuit breaker is open
	assert.Equal(t, http.StatusServiceUnavailable, *res.Code)
	assert.Equal(t, "network error, unable to reach weather service", *res.Message)
}
