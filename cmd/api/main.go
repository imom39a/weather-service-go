package main

import (
	"net/http"
	"weathe-service/internal/api"
	handlers "weathe-service/internal/api/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Up and running!")
	})

	api.RegisterHandlers(e, handlers.NewCompositeHandler())
	e.Logger.Fatal(e.Start(":8080"))
}
