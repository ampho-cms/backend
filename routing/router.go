// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package routing

import "net/http"

type Router interface {
	// Vars returns request variables.
	Vars(req *http.Request) map[string]string

	// Handle registers a new route.
	Handle(path string, handler RequestHandler) Route

	// AddMiddleware appends a middleware to the chain.
	AddMiddleware(handler MiddlewareHandler)

	// ServeHTTP dispatches the handler registered in the matched route.
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
