// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package routing

import (
	"net/http"
)

// Request is the http.Request wrapper with several helper methods.
type Request struct {
	request *http.Request
	router  Router
}

// RequestHandler is the HTTP request handler.
type RequestHandler func(req *Request, resp *Response)

// MiddlewareHandler is the middleware HTTP request handler.
// Should return false if the next middleware in chain must not be called.
type MiddlewareHandler func(req *Request, resp *Response) bool

// Request returns the underlying http.Request object.
func (r *Request) Request() *http.Request {
	return r.request
}

// Vars returns request variables.
func (r *Request) Vars(req *http.Request) map[string]string {
	return r.router.Vars(req)
}

// Var returns a request variable value.
func (r *Request) Var(k string) string {
	return r.router.Vars(r.Request())[k]
}
