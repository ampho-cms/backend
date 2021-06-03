// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"ampho.xyz/ampho/config"
	"ampho.xyz/ampho/logging"
	"ampho.xyz/ampho/routing"
	"ampho.xyz/ampho/server"
)

// Service is the service interface.
type Service interface {
	Name() string
	Signature() string
	Mode() string
	Cfg() config.Config
	Log() logging.Logger
	Router() routing.Router
	Server() *http.Server
	Start()
	Stop()
	Run()
}

// Base is the base service structure.
type Base struct {
	name   string
	mode   string
	cfg    config.Config
	log    logging.Logger
	router routing.Router
	srv    *http.Server
}

// Name returns service name.
func (s *Base) Name() string {
	return s.name
}

// Signature returns signature used as a 'Server' HTTP header.
func (s *Base) Signature() string {
	return s.name
}

// Mode returns service mode.
func (s *Base) Mode() string {
	return s.mode
}

// Cfg returns service configuration engine.
func (s *Base) Cfg() config.Config {
	return s.cfg
}

// Log returns service logger.
func (s *Base) Log() logging.Logger {
	return s.log
}

// Router returns service router.
func (s *Base) Router() routing.Router {
	return s.router
}

// Server returns service HTTP server.
func (s *Base) Server() *http.Server {
	return s.srv
}

// Start starts the service. Assumed to be run as a goroutine.
func (s *Base) Start() {
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		s.log.DebugF("%s has been stopped", s.name)
	} else {
		s.log.Error(err)
	}
}

// Stop stops the service.
func (s *Base) Stop() {
	s.log.Debug("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GetDuration("service.shutdownTimeout"))
	defer cancel()

	// Ask server to gracefully shutdown. After it finish shutting down
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Error(err.Error())
	}

	// Flush logs
	_ = s.log.Sync()
}

// Run starts the server and waits for SIGINT. Usually it should be a last call in the `main()`.
func (s *Base) Run() {
	go s.Start()
	s.log.DebugF("%s started at http://%s in %s mode", s.name, s.srv.Addr, s.mode)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()
}

// New creates a new Base instance.
func New(name, mode string, cfg config.Config, log logging.Logger, router routing.Router, srv *http.Server) *Base {
	// Sanitize mode
	switch mode {
	case ModeDevelopment, ModeProduction:
	default:
		mode = ModeDevelopment
	}

	return &Base{name, mode, cfg, log, router, srv}
}

// ConfigDefaults returns service configuration defaults.
func ConfigDefaults() map[string]interface{} {
	return map[string]interface{}{
		"service.mode":            ModeDevelopment,
		"service.shutdownTimeout": DftShutdownTimeout,
		"logging.console.enabled": DftLoggingConsoleEnabled,
		"logging.file.size":       DftLoggingFileSize,
		"logging.file.age":        DftLoggingFileAge,
		"logging.file.backups":    DftLoggingFileBackups,
		"network.address":         DftNetAddr,
		"network.readTimeout":     DftNetReadTimeout,
		"network.writeTimeout":    DftNetWriteTimeout,
	}
}

// NewDefault creates a new base service instance using default configuration.
func NewDefault(name string) (*Base, error) {
	// Configuration engine
	cfg, err := config.NewDefault(name, ConfigDefaults())
	if err != nil {
		return nil, err
	}

	// Operating mode
	mode := cfg.GetString("service.mode")
	switch mode {
	case ModeDevelopment, ModeProduction:
	default:
		mode = ModeDevelopment
	}

	// Logging
	var log logging.Logger
	logPath := cfg.GetString("logging.file.path")
	if logPath != "" {
		logPath = strings.ReplaceAll(logPath, "$SERVICE_NAME", name)
		log, err = logging.NewDefaultRotatingFile(
			mode,
			logPath,
			cfg.GetInt("logging.file.size"),
			cfg.GetInt("logging.file.age"),
			cfg.GetInt("logging.file.backups"),
			cfg.GetBool("logging.console.enabled"),
		)
		if err != nil {
			return nil, err
		}
	} else if cfg.GetBool("logging.console.enabled") {
		if log, err = logging.NewDefaultConsole(mode); err != nil {
			return nil, err
		}
	}

	// Router, server and the service
	router := routing.NewDefault()
	srv := server.NewDefault(
		cfg.GetString("network.address"),
		router,
		cfg.GetDuration("network.readTimeout"),
		cfg.GetDuration("network.writeTimeout"),
	)
	svc := New(name, mode, cfg, log, router, srv)

	// Middlewares
	router.AddMiddleware(svc.ServerSignatureMiddleware)
	if mode == ModeDevelopment {
		router.AddMiddleware(svc.RequestLogDebugMiddleware)
	}

	return svc, nil
}

// NewTesting creates a new base service instance suitable for unit tests.
func NewTesting(name string) *Base {
	// In-memory config
	cfg := config.NewMemory()
	for k, v := range ConfigDefaults() {
		cfg.SetDefault(k, v)
	}

	// In-memory logging
	log := logging.NewMemory(logging.LDebug)

	router := routing.NewDefault()

	srv := server.NewDefault(
		cfg.GetString("network.address"),
		router,
		cfg.GetDuration("network.readTimeout"),
		cfg.GetDuration("network.writeTimeout"),
	)

	return New(name, ModeDevelopment, cfg, log, router, srv)
}
