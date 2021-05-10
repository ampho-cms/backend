// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package routing

import (
	"github.com/gorilla/mux"
	"net/http"
)

// GorillaMuxRouter is GorillaMux backend.
type GorillaMuxRouter struct {
	backend *mux.Router
}

// GorillaMuxRoute is the wrapper of GorillaMux route.
type GorillaMuxRoute struct {
	backend *mux.Route
}

// Vars returns request variables.
func (r *GorillaMuxRouter) Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// Handle registers a new route.
func (r *GorillaMuxRouter) Handle(path string, handler RequestHandler) Route {
	muxRoute := r.backend.HandleFunc(path, func(writer http.ResponseWriter, req *http.Request) {
		handler(&Request{req, r}, &Response{writer})
	})

	return &GorillaMuxRoute{muxRoute}
}

// AddMiddleware appends a middleware to the chain.
func (r *GorillaMuxRouter) AddMiddleware(handler MiddlewareHandler) {
	r.backend.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if handler(&Request{request: req}, &Response{writer: resp}) {
				next.ServeHTTP(resp, req)
			}
		})
	})
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *GorillaMuxRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.backend.ServeHTTP(writer, request)
}

// NewGorillaMux creates a new gorilla/mux backend router.
func NewGorillaMux() *GorillaMuxRouter {
	return &GorillaMuxRouter{mux.NewRouter()}
}
