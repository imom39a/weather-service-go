// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

// GetWeatherParams defines parameters for GetWeather.
type GetWeatherParams struct {
	// Lat Latitude
	Lat float32 `form:"lat" json:"lat"`

	// Lon Longitude
	Lon float32 `form:"lon" json:"lon"`
}