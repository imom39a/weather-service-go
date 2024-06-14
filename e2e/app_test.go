package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"weathe-service/internal/api"
	handlers "weathe-service/internal/api/handler"

	"github.com/labstack/echo/v4"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.prepare()
			rec := httptest.NewRecorder()
			e := echo.New()
			api.RegisterHandlers(e, handlers.NewCompositeHandler())
			e.ServeHTTP(rec, req)

			body, _ := io.ReadAll(rec.Body)
			tt.assertFunc(t, rec.Result(), body)
		})
	}
}
