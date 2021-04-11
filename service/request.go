// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RequestHandler is the service HTTP request function
type RequestHandler func(*Service, *Request, *Response)

// Request is the service HTTP request.
type Request struct {
	request *http.Request
}

// Request returns an underlying http.Request object.
func (r *Request) Request() *http.Request {
	return r.request
}

// Vars returns request variables.
func (r *Request) Vars() map[string]string {
	return mux.Vars(r.request)
}

// Var returns a request variable value.
func (r *Request) Var(k string) string {
	return r.Vars()[k]
}
