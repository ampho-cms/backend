// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "time"

const (
	DftShutdownTimeout = time.Second * 15 // service shutdown timeout

	DftLoggingLevel = "debug" // logging level

	DftLoggingFileSize    = 100  // logging file size
	DftLoggingFileAge     = 30   // logging file age
	DftLoggingFileBackups = 12   // logging backup files number
	DftLoggingFileStdout  = true // whether to duplicate file logging to stdout

	DftAddress      = "127.0.0.1:8765" // HTTP server address
	DftReadTimeout  = time.Second * 15 // network read timeout
	DftWriteTimeout = time.Second * 15 // network write timeout

	DftSecuritySigningMethod = "HS256" // security singing method
)
