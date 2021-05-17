// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "time"

const (
	ModeDevelopment = "development" // development mode
	ModeProduction  = "production"  // production mode

	DftShutdownTimeout = time.Second * 15 // default service shutdown timeout

	DftLoggingConsoleEnabled = true // default console logging state
	DftLoggingFileSize       = 100  // default logging file size
	DftLoggingFileAge        = 30   // default logging file age
	DftLoggingFileBackups    = 12   // default logging backup files number

	DftNetAddr         = "127.0.0.1:8765" // default HTTP server address
	DftNetReadTimeout  = time.Second * 15 // default network read timeout
	DftNetWriteTimeout = time.Second * 15 // default network write timeout
)
