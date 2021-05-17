// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"ampho.xyz/ampho/config"
	"ampho.xyz/ampho/logging"
	"ampho.xyz/ampho/routing"
	"github.com/spf13/viper"
)

// Service is the base service structure.
type Service struct {
	name   string
	mode   string
	cfg    config.Config
	log    logging.Logger
	router routing.Router
	server *http.Server
}

// Name returns service name.
func (s *Service) Name() string {
	return s.name
}

// Signature returns signature used as a 'Server' HTTP header.
func (s *Service) Signature() string {
	return s.name
}

// Mode returns service mode.
func (s *Service) Mode() string {
	return s.mode
}

// Cfg returns the configurator.
func (s *Service) Cfg() config.Config {
	return s.cfg
}

// Log returns the logger.
func (s *Service) Log() logging.Logger {
	return s.log
}

// Router returns the router.
func (s *Service) Router() routing.Router {
	return s.router
}

// Server returns the server.
func (s *Service) Server() *http.Server {
	return s.server
}

// Start starts the service. Intended to be run in a goroutine.
func (s *Service) Start() {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		s.log.DebugF("%s has been stopped", s.name)
	} else {
		s.log.Error(err)
	}
}

// Stop stops the service.
func (s *Service) Stop() {
	s.log.Debug("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GetDuration("service.shutdownTimeout"))
	defer cancel()

	// Ask server to gracefully shutdown. After it finish shutting down
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error(err.Error())
	}

	// Flush logs
	_ = s.log.Sync()
}

// Run starts the server and waits for SIGINT. Usually it should be a last call in the `main()`.
func (s *Service) Run() {
	go s.Start()
	s.log.DebugF("%s started at http://%s in %s mode", s.name, s.server.Addr, s.mode)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()
}

// NewDefaultConfig creates and configures a new default configurator.
func NewDefaultConfig(name string) (config.Config, error) {
	// Executable directory path
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	// Configuration engine
	vp := viper.New()
	vp.SetConfigName(name)
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME/." + name)
	vp.AddConfigPath(execDir)
	vp.AddConfigPath(".")
	cfg := config.NewViper(vp)

	// Configuration defaults
	cfg.SetDefault("service.mode", ModeDevelopment)
	cfg.SetDefault("service.shutdownTimeout", DftShutdownTimeout)
	cfg.SetDefault("logging.console.enabled", DftLoggingConsoleEnabled)
	cfg.SetDefault("logging.file.size", DftLoggingFileSize)
	cfg.SetDefault("logging.file.age", DftLoggingFileAge)
	cfg.SetDefault("logging.file.backups", DftLoggingFileBackups)
	cfg.SetDefault("network.address", DftNetAddr)
	cfg.SetDefault("network.readTimeout", DftNetReadTimeout)
	cfg.SetDefault("network.writeTimeout", DftNetWriteTimeout)

	// Load configuration
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return cfg, nil
}

// NewDefaultRouter creates a new default router.
func NewDefaultRouter() routing.Router {
	return routing.NewGorillaMux()
}

// NewDefaultServer instantiates a new HTTP server using default configuration.
func NewDefaultServer(cfg config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         cfg.GetString("network.address"),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}
}

// New creates a new Service instance.
func New(name, mode string, cfg config.Config, log logging.Logger, router routing.Router, srv *http.Server) *Service {
	// Sanitize mode
	switch mode {
	case ModeDevelopment, ModeProduction:
	default:
		mode = ModeDevelopment
	}

	return &Service{name, mode, cfg, log, router, srv}
}

// NewDefault creates a new Service instance using default configuration.
func NewDefault(name string) (*Service, error) {
	// Configuration
	cfg, err := NewDefaultConfig(name)
	if err != nil {
		return nil, err
	}

	// Determine mode
	mode := cfg.GetString("service.mode")
	switch mode {
	case ModeDevelopment, ModeProduction:
	default:
		mode = ModeDevelopment
	}

	// Logger
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

	// Service
	router := NewDefaultRouter()
	svc := New(name, mode, cfg, log, router, NewDefaultServer(cfg, router))

	// Middlewares
	router.AddMiddleware(svc.ServerSignatureMiddleware)
	if mode == ModeDevelopment {
		router.AddMiddleware(svc.RequestLogDebugMiddleware)
	}

	return svc, nil
}

// NewTesting creates a new Service instance suitable for using in unit tests.
func NewTesting(name string) *Service {
	log := logging.NewMemory(logging.LDebug)
	cfg := config.NewMemory()
	router := NewDefaultRouter()
	srv := NewDefaultServer(cfg, router)

	return New(name, ModeDevelopment, cfg, log, router, srv)
}
