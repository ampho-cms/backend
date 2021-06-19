// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zapLogger is the Zap logger backend.
type zapLogger struct {
	backend *zap.SugaredLogger
	level   Level
}

// Backend returns logger backend.
func (z *zapLogger) Backend() interface{} {
	return z.backend
}

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
func (z *zapLogger) Sync() error {
	return z.backend.Sync()
}

// Level returns logger level
func (z *zapLogger) Level() Level {
	return z.level
}

// Debug logs a message at debug level.
func (z *zapLogger) Debug(args ...interface{}) {
	z.backend.Debug(args...)
}

// DebugF logs a formatted message at debug level.
func (z *zapLogger) DebugF(tpl string, args ...interface{}) {
	z.backend.Debugf(tpl, args...)
}

// Info logs a message at info level.
func (z *zapLogger) Info(args ...interface{}) {
	z.backend.Info(args...)
}

// InfoF logs a formatted message at info level.
func (z *zapLogger) InfoF(tpl string, args ...interface{}) {
	z.backend.Infof(tpl, args...)
}

// Warn logs a message at warn level.
func (z *zapLogger) Warn(args ...interface{}) {
	z.backend.Warn(args...)
}

// WarnF logs a formatted message at warning level.
func (z *zapLogger) WarnF(tpl string, args ...interface{}) {
	z.backend.Warnf(tpl, args...)
}

// Error logs a message at error level.
func (z *zapLogger) Error(args ...interface{}) {
	z.backend.Error(args...)
}

// ErrorF logs a formatted message at error level.
func (z *zapLogger) ErrorF(tpl string, args ...interface{}) {
	z.backend.Errorf(tpl, args...)
}

// Panic logs a message at panic level and then panics.
func (z *zapLogger) Panic(args ...interface{}) {
	z.backend.Panic(args...)
}

// PanicF logs a formatted message at panic level and then panics.
func (z *zapLogger) PanicF(tpl string, args ...interface{}) {
	z.backend.Panicf(tpl, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func (z *zapLogger) Fatal(args ...interface{}) {
	z.backend.Fatal(args...)
}

// FatalF logs a formatted message at fatal level and then calls os.Exit(1).
func (z *zapLogger) FatalF(tpl string, args ...interface{}) {
	z.backend.Fatalf(tpl, args...)
}

// NewZap creates a new zapLogger logger
func NewZap(backend *zap.Logger, level Level) *zapLogger {
	return &zapLogger{backend.Sugar(), level}
}

// getZapLevel returns zapLogger logging level corresponding to Ampho's one.
func getZapLevel(amohoLevel Level) zapcore.LevelEnabler {
	switch amohoLevel {
	case LInfo:
		return zapcore.InfoLevel
	case LWarn:
		return zapcore.WarnLevel
	case LError:
		return zapcore.ErrorLevel
	case LPanic:
		return zapcore.PanicLevel
	case LFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// newZapConsole creates a zapLogger console logger.
func newZapConsole(level Level) (Logger, error) {
	w, _, err := zap.Open("stdout")
	if err != nil {
		return nil, err
	}

	var encCfg zapcore.EncoderConfig
	if level == LDebug {
		encCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encCfg = zap.NewProductionEncoderConfig()
	}

	zapLogger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), w, getZapLevel(level)))
	return NewZap(zapLogger, level), nil
}

// newZapRotatingFile creates a zapLogger file rolling logger.
func newZapRotatingFile(level Level, path string, size, age, backups int, stdout bool) (Logger, error) {
	// Sanitize relative file path
	if !filepath.IsAbs(path) {
		execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}
		path = filepath.Join(execDir, path)
	}

	var writers []zapcore.WriteSyncer

	if stdout {
		w, _, err := zap.Open("stdout")
		if err != nil {
			return nil, err
		}
		writers = append(writers, w)
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    size, // megabytes
		MaxAge:     age,  // days
		MaxBackups: backups,
	})
	writers = append(writers, w)

	var encCfg zapcore.EncoderConfig
	if level == LDebug {
		encCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encCfg = zap.NewProductionEncoderConfig()
	}

	w = zapcore.NewMultiWriteSyncer(writers...)

	zapLogger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), w, getZapLevel(level)))

	return NewZap(zapLogger, level), nil
}
