// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package logging

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewDefaultConsole creates a default console logger.
func NewDefaultConsole(mode string) (Logger, error) {
	var (
		encCfg   zapcore.EncoderConfig
		zapLevel zapcore.LevelEnabler
	)

	w, _, err := zap.Open("stdout")
	if err != nil {
		return nil, err
	}

	if mode == ModeProduction {
		encCfg = zap.NewProductionEncoderConfig()
		zapLevel = zapcore.InfoLevel
	} else {
		encCfg = zap.NewDevelopmentEncoderConfig()
		zapLevel = zapcore.DebugLevel
	}

	log.SetOutput(w)

	return NewZap(zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), w, zapLevel))), nil
}

// NewDefaultRotatingFile creates a default file rolling logger.
func NewDefaultRotatingFile(mode string, path string, size, age, backups int, stdout bool) (Logger, error) {
	var (
		encCfg   zapcore.EncoderConfig
		zapLevel zapcore.LevelEnabler
	)

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

	if mode == ModeProduction {
		encCfg = zap.NewProductionEncoderConfig()
		zapLevel = zapcore.InfoLevel
	} else {
		encCfg = zap.NewDevelopmentEncoderConfig()
		zapLevel = zapcore.DebugLevel
	}

	w = zapcore.NewMultiWriteSyncer(writers...)
	log.SetOutput(w)

	return NewZap(zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), w, zapLevel))), nil
}
