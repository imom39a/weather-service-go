package server

import (
	"net/http"
	errors "weathe-service/common/error"
	"weathe-service/internal/api"

	"github.com/labstack/echo/v4"
)

func CreateServer() *echo.Echo {
	e := echo.New()
	e.Use(errorHandlerMiddleware)
	return e
}

func errorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		gerr := next(c)
		if gerr != nil {
			switch e := gerr.(type) {
			case *errors.APIError:
				return createErrorResponse(c, e.HTTPErrCode, e.Description)
			case *echo.HTTPError:
				if msg, ok := e.Message.(string); ok {
					return createErrorResponse(c, e.Code, msg)
				}
				return createErrorResponse(c, e.Code, "HTTP error occurred")
			default:
				return createErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
			}
		}
		return nil
	}
}

func createErrorResponse(c echo.Context, statusCode int, description string) error {
	return c.JSON(statusCode, &api.ErrorResponse{
		Code:    &statusCode,
		Message: &description,
	})
}
