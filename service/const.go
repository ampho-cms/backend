// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "time"

const (
	ModeDevelopment = "development" // development mode
	ModeProduction  = "production"  // production mode

	DftNetAddr         = "127.0.0.1:8765" // default HTTP server address
	DftShutdownTimeout = time.Second * 15 // default service shutdown timeout
	DftNetReadTimeout  = time.Second * 15 // default network read timeout
	DftNetWriteTimeout = time.Second * 15 // default network write timeout
)
