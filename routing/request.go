// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package routing provides routing related things.
package routing

import (
	"net/http"
)

// RequestHandler is the HTTP request handler structure.
type RequestHandler func(*Request, *Response)

// Request is the service HTTP request structure.
type Request struct {
	request *http.Request
	router  Router
}

// Request returns the underlying http.Request object.
func (r *Request) Request() *http.Request {
	return r.request
}

// Vars returns request variables. Shortcut method.
func (r *Request) Vars(req *http.Request, k string) map[string]string {
	return r.router.Vars(req)
}

// Var returns a request variable value. Shortcut method.
func (r *Request) Var(req *http.Request, k string) string {
	return r.router.Vars(req)[k]
}
