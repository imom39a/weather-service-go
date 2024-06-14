package api_test

import (
	"net/http"
	"testing"
)

type ApiTest struct {
	name        string
	description string
	prepare     func() *http.Request
	assertFunc  func(t *testing.T, rec *http.Response, body []byte)
}
