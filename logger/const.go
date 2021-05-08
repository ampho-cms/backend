// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logger

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// LDebug logs are typically voluminous, and are usually disabled in production.
	LDebug Level = iota - 1

	// LInfo is the default logging priority.
	LInfo

	// LWarn logs are more important than Info, but don't need individual human review.
	LWarn

	// LError logs are high-priority. If an app is running smoothly, it shouldn't generate any error-level logs.
	LError

	// LPanic logs a message, then panics.
	LPanic

	// LFatal logs a message, then calls os.Exit(1).
	LFatal

	minLevel = LDebug
	maxLevel = LFatal
)
