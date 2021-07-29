// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package servicetest provides helper functions for testing services.
package servicetest

import (
	"ampho.xyz/service"
	"net/http"
	"net/http/httptest"
)

// DoRequest performs a request to the service and writes a response.
func DoRequest(svc service.Service, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	svc.Router().ServeHTTP(rr, req)
	return rr
}
