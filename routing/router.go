// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package routing provides routing related things.
package routing

import "net/http"

type Router interface {
	// Vars returns request variables.
	Vars(req *http.Request) map[string]string

	// AddHandler registers a new route.
	AddHandler(path string, handler RequestHandler) Route

	// AddMiddleware appends a middleware to the chain.
	AddMiddleware(handler RequestHandler)

	// ServeHTTP dispatches the handler registered in the matched route.
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
