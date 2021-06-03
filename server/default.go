// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package server provides primitives, structures and functions related to HTTP server.
package server

import (
	"net/http"
	"time"
)

// NewDefault creates a default HTTP server.
func NewDefault(addr string, handler http.Handler, rTOut, wTOut time.Duration) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         addr,
		ReadTimeout:  rTOut,
		WriteTimeout: wTOut,
	}
}
