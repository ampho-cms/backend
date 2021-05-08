// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "time"

const (
	DftNetAddr         = "127.0.0.1:8765"
	DftNetReadTimeout  = 15
	DftNetWriteTimeout = 15

	DftShutdownTimeout = 15 * time.Second

	ModeDevelopment = "development"
	ModeProduction  = "production"
)
