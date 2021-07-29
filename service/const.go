// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "time"

const (
	DftShutdownTimeout = time.Second * 15 // service shutdown timeout



	DftAddress      = "127.0.0.1:8765" // HTTP server address
	DftReadTimeout  = time.Second * 15 // network read timeout
	DftWriteTimeout = time.Second * 15 // network write timeout
)
