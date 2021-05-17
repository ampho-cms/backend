// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logging

import "go.uber.org/zap"

// Zap is the Zap logger backend.
type Zap struct {
	backend *zap.SugaredLogger
}

// Backend returns logger backend.
func (z *Zap) Backend() interface{} {
	return z.backend
}

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
func (z *Zap) Sync() error {
	return z.backend.Sync()
}

// Debug logs a message at debug level.
func (z *Zap) Debug(args ...interface{}) {
	z.backend.Debug(args...)
}

// DebugF logs a formatted message at debug level.
func (z *Zap) DebugF(tpl string, args ...interface{}) {
	z.backend.Debugf(tpl, args...)
}

// Info logs a message at info level.
func (z *Zap) Info(args ...interface{}) {
	z.backend.Info(args...)
}

// InfoF logs a formatted message at info level.
func (z *Zap) InfoF(tpl string, args ...interface{}) {
	z.backend.Infof(tpl, args...)
}

// Warn logs a message at warn level.
func (z *Zap) Warn(args ...interface{}) {
	z.backend.Warn(args...)
}

// WarnF logs a formatted message at warning level.
func (z *Zap) WarnF(tpl string, args ...interface{}) {
	z.backend.Warnf(tpl, args...)
}

// Error logs a message at error level.
func (z *Zap) Error(args ...interface{}) {
	z.backend.Error(args...)
}

// ErrorF logs a formatted message at error level.
func (z *Zap) ErrorF(tpl string, args ...interface{}) {
	z.backend.Errorf(tpl, args...)
}

// Panic logs a message at panic level and then panics.
func (z *Zap) Panic(args ...interface{}) {
	z.backend.Panic(args...)
}

// PanicF logs a formatted message at panic level and then panics.
func (z *Zap) PanicF(tpl string, args ...interface{}) {
	z.backend.Panicf(tpl, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func (z *Zap) Fatal(args ...interface{}) {
	z.backend.Fatal(args...)
}

// FatalF logs a formatted message at fatal level and then calls os.Exit(1).
func (z *Zap) FatalF(tpl string, args ...interface{}) {
	z.backend.Fatalf(tpl, args...)
}

// NewZap creates a new Zap logger
func NewZap(backend *zap.Logger) *Zap {
	return &Zap{backend.Sugar()}
}
