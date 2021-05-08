// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package routing provides routing related things.
package routing

import (
	"github.com/gorilla/mux"
	"net/http"

	"ampho.xyz/ampho/logger"
)

type GorillaMuxRouter struct {
	backend *mux.Router
	log     logger.Logger
}

type GorillaMuxRoute struct {
	backend *mux.Route
}

// Vars returns request variables.
func (r *GorillaMuxRouter) Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// AddHandler registers a new route.
func (r *GorillaMuxRouter) AddHandler(path string, handler RequestHandler) Route {
	muxRoute := r.backend.HandleFunc(path, func(writer http.ResponseWriter, req *http.Request) {
		handler(&Request{req, r}, &Response{writer})
	})

	r.log.DebugF("route registered: %s", path)

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

	r.log.DebugF("middleware registered: %v", handler)
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *GorillaMuxRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.backend.ServeHTTP(writer, request)
}

// NewGorillaMux creates a new gorilla/mux backend router.
func NewGorillaMux(log logger.Logger) Router {
	return &GorillaMuxRouter{
		mux.NewRouter(),
		log,
	}
}
