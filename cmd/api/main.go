package main

import (
	"net/http"
	"weathe-service/common/logger"
	"weathe-service/common/server"
	"weathe-service/internal/api"
	handlers "weathe-service/internal/api/handler"
)

func main() {
	e := server.CreateServer()
	e.Use(server.GetSwaggerValidatorMiddleware("spec/weather-service.yaml"))
	api.RegisterHandlers(e, handlers.NewCompositeHandler(logger.NewLogger(), getHttpClients()))
	e.Logger.Fatal(e.Start(":8080"))
}

func getHttpClients() *http.Client {
	client := &http.Client{}
	return client
}
