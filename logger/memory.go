// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Memory is the in memory logger backend.
type Memory struct {
	level   Level
	count   int64
	storage []string
}

// log logs a message
func (m *Memory) log(level Level, tpl string, args ...interface{}) string {
	if level < m.level {
		return ""
	}

	lName := ""
	switch level {
	case LDebug:
		lName = "DEBUG"
	case LInfo:
		lName = "INFO"
	case LWarn:
		lName = "WARN"
	case LError:
		lName = "ERROR"
	case LPanic:
		lName = "PANIC"
	case LFatal:
		lName = "FATAL"
	default:
		lName = "UNKNOWN"
	}

	msg := ""
	if tpl != "" {
		msg = fmt.Sprintf(tpl, args...)
	} else {
		for _, v := range args {
			msg = fmt.Sprintf("%s %v", msg, v)
		}
	}

	msg = fmt.Sprintf("%v\t%s\t%s", time.Now(), lName, msg)
	m.storage = append(m.storage, msg)
	log.Printf(msg)

	m.count++

	return msg
}

// Backend returns logger backend.
func (m *Memory) Backend() interface{} {
	return m.storage
}

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
func (m *Memory) Sync() error {
	return nil
}

// Debug logs a message at debug level.
func (m *Memory) Debug(args ...interface{}) {
	m.log(LDebug, "", args...)
}

// DebugF logs a formatted message at debug level.
func (m *Memory) DebugF(tpl string, args ...interface{}) {
	m.log(LDebug, tpl, args...)
}

// Info logs a message at info level.
func (m *Memory) Info(args ...interface{}) {
	m.log(LInfo, "", args...)
}

// InfoF logs a formatted message at info level.
func (m *Memory) InfoF(tpl string, args ...interface{}) {
	m.log(LInfo, tpl, args...)
}

// Warn logs a message at warn level.
func (m *Memory) Warn(args ...interface{}) {
	m.log(LWarn, "", args...)
}

// WarnF logs a formatted message at warning level.
func (m *Memory) WarnF(tpl string, args ...interface{}) {
	m.log(LWarn, tpl, args...)
}

// Error logs a message at error level.
func (m *Memory) Error(args ...interface{}) {
	m.log(LError, "", args...)
}

// ErrorF logs a formatted message at error level.
func (m *Memory) ErrorF(tpl string, args ...interface{}) {
	m.log(LError, tpl, args...)
}

// Panic logs a message at panic level and then panics.
func (m *Memory) Panic(args ...interface{}) {
	if msg := m.log(LPanic, "", args...); msg != "" {
		panic(msg)
	}
}

// PanicF logs a formatted message at panic level and then panics.
func (m *Memory) PanicF(tpl string, args ...interface{}) {
	panic(m.log(LPanic, tpl, args...))
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func (m *Memory) Fatal(args ...interface{}) {
	if m.log(LFatal, "", args...) != "" {
		os.Exit(1)
	}
}

// FatalF logs a formatted message at fatal level and then calls os.Exit(1).
func (m *Memory) FatalF(tpl string, args ...interface{}) {
	if m.log(LFatal, tpl, args...) != "" {
		os.Exit(1)
	}
}

// NewMemory creates a new Memory logger.
func NewMemory(level Level) *Memory {
	return &Memory{level, 0, make([]string, 0)}
}
