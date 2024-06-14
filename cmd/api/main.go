package main

import (
	"weathe-service/common/logger"
	"weathe-service/common/server"
	"weathe-service/internal/api"
	handlers "weathe-service/internal/api/handler"
)

func main() {
	e := server.CreateServer()
	api.RegisterHandlers(e, handlers.NewCompositeHandler(logger.NewLogger()))
	e.Logger.Fatal(e.Start(":8080"))
}
