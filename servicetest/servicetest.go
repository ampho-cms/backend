// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package servicetest provides helper functions for testing services.
package servicetest

import (
	"net/http"
	"net/http/httptest"

	"ampho.xyz/ampho/service"
)

// DoRequest performs a request to the service and writes a response.
func DoRequest(svc service.Service, method, target string) *http.Response {
	req := httptest.NewRequest(method, target, nil)
	respW := httptest.NewRecorder()
	svc.Server().Handler.ServeHTTP(respW, req)

	return respW.Result()
}
