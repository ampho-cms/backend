// Package logger contains logging related things
package logging

// Logger is the logger interface
type Logger interface {
	// Backend returns logger backend.
	Backend() interface{}

	// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
	Sync() error

	// Debug logs a message at debug level.
	Debug(args ...interface{})

	// DebugF logs a formatted message at debug level.
	DebugF(tpl string, args ...interface{})

	// Info logs a message at info level.
	Info(args ...interface{})

	// InfoF logs a formatted message at info level.
	InfoF(tpl string, args ...interface{})

	// Warn logs a message at warn level.
	Warn(args ...interface{})

	// WarnF logs a formatted message at warning level.
	WarnF(tpl string, args ...interface{})

	// Error logs a message at error level.
	Error(args ...interface{})

	// ErrorF logs a formatted message at error level.
	ErrorF(tpl string, args ...interface{})

	// Panic logs a message at panic level and then panics.
	Panic(args ...interface{})

	// PanicF logs a formatted message at panic level and then panics.
	PanicF(tpl string, args ...interface{})

	// Fatal logs a message at fatal level and then calls os.Exit(1).
	Fatal(args ...interface{})

	// FatalF logs a formatted message at fatal level and then calls os.Exit(1).
	FatalF(tpl string, args ...interface{})
}
