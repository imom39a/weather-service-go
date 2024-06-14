package server

import (
	"net/http"
	"os"
	errors "weathe-service/common/error"
	"weathe-service/internal/api"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapi_middleware "github.com/oapi-codegen/echo-middleware"
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

func GetSwaggerValidatorMiddleware(filePath string) echo.MiddlewareFunc {
	swaggerLoader := openapi3.NewLoader()
	swagger, err := swaggerLoader.LoadFromFile(filePath)
	if err != nil {
		os.Exit(1)
	}

	swagger.Servers = make(openapi3.Servers, 0)
	swagger.Servers = append(swagger.Servers, &openapi3.Server{
		URL: "",
	})

	options := &oapi_middleware.Options{
		SilenceServersWarning: true,
		Options: openapi3filter.Options{
			IncludeResponseStatus: true,
			SkipSettingDefaults:   true,
		},
	}

	return oapi_middleware.OapiRequestValidatorWithOptions(swagger, options)
}
