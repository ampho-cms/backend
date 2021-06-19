// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logging

import (
	"strings"
)

var defaultLogger Logger

var levelsByName = map[string]Level{
	"debug":   LDebug,
	"info":    LInfo,
	"warn":    LWarn,
	"warning": LWarn,
	"err":     LError,
	"error":   LError,
	"panic":   LPanic,
	"fatal":   LFatal,
}

// SetLogger sets a default logger
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// GetLogger returns a default logger
func GetLogger() Logger {
	if defaultLogger == nil {
		panic("Default logger is not set. Did you forget to call SetLogger()?")
	}

	return defaultLogger
}

// GetLevelByName return a logging level by its name.
func GetLevelByName(name string) Level {
	lv, ok := levelsByName[strings.ToLower(name)]
	if !ok {
		lv = LDebug
	}

	return lv
}

// NewConsoleLogger creates a new console logger using default configuration.
func NewConsoleLogger(level Level) (Logger, error) {
	return newZapConsole(level)
}

// NewFileLogger creates a new file logger using default configuration.
func NewFileLogger(path string, size, age, backups int, level Level, stdout bool) (Logger, error) {
	return newZapRotatingFile(level, path, size, age, backups, stdout)
}

// Debug logs a message at debug level.
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// DebugF logs a formatted message at debug level.
func DebugF(tpl string, args ...interface{}) {
	GetLogger().DebugF(tpl, args...)
}

// Info logs a message at info level.
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// InfoF logs a formatted message at info level.
func InfoF(tpl string, args ...interface{}) {
	GetLogger().InfoF(tpl, args...)
}

// Warn logs a message at warn level.
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// WarnF logs a formatted message at warning level.
func WarnF(tpl string, args ...interface{}) {
	GetLogger().WarnF(tpl, args...)
}

// Error logs a message at error level.
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// ErrorF logs a formatted message at error level.
func ErrorF(tpl string, args ...interface{}) {
	GetLogger().ErrorF(tpl, args...)
}

// Panic logs a message at panic level and then panics.
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// PanicF logs a formatted message at panic level and then panics.
func PanicF(tpl string, args ...interface{}) {
	GetLogger().PanicF(tpl, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// FatalF logs a formatted message at fatal level and then calls os.Exit(1).
func FatalF(tpl string, args ...interface{}) {
	GetLogger().FatalF(tpl, args...)
}
