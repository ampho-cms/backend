// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ampho-cms/backend/config"
	"github.com/ampho-cms/backend/logger"
	"github.com/ampho-cms/backend/routing"
)

// Service is the service.
type Service struct {
	name    string
	version string
	mode    string
	cfg     config.Config
	log     logger.Logger
	router  routing.Router
	server  *http.Server
}

// Name returns service name.
func (s *Service) Name() string {
	return s.name
}

// Signature returns service signature used as a 'Server' HTTP header.
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
func (s *Service) Log() logger.Logger {
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

// Start starts the service in blocking mode.
// This method is usually should be run in a goroutine.
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
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GetDuration("shutdownTimeout"))
	defer cancel()

	// Ask server to gracefully shutdown. After it finish shutting down
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error(err.Error())
	}

	_ = s.log.Sync()
}

// Run starts the server and waits for SIGINT.
//
// Usually this function should be a last call in your `main()`.
func (s *Service) Run() {
	go s.Start()
	s.log.DebugF("%s started at http://%s in %s mode", s.name, s.server.Addr, s.mode)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()
}

// NewDefaultLogger and configures a new default logger.
func NewDefaultLogger(mode string) logger.Logger {
	var (
		err   error
		zpCfg zap.Config
		zp    *zap.Logger
	)

	if mode == ModeProduction {
		zpCfg = zap.NewProductionConfig()
	} else {
		zpCfg = zap.NewDevelopmentConfig()
	}

	//zpCfg.DisableCaller = true // disable caller logging since it doesn't work correctly
	zp, err = zpCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Errorf("Fatal error while initializing a logger: %s\n", err))
	}
	l := logger.NewZap(zp)
	l.DebugF("logging initialized in %s mode", mode)

	return l
}

// NewDefaultConfig creates and configures a new default configurator.
func NewDefaultConfig(name string) config.Config {
	// Configuration engine
	vp := viper.New()
	vp.SetConfigName("." + name + ".yaml")
	vp.AddConfigPath(".")
	cfg := config.NewViper(vp)

	// Configuration defaults
	cfg.SetDefault("mode", ModeDevelopment)
	cfg.SetDefault("address", DftNetAddr)
	cfg.SetDefault("readTimeout", DftNetReadTimeout)
	cfg.SetDefault("writeTimeout", DftNetWriteTimeout)
	cfg.SetDefault("shutdownTimeout", DftShutdownTimeout)

	// Load configuration
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("Fatal error while loading config file: %s \n", err))
		}
	}

	return cfg
}

// NewDefaultRouter creates and configures a default router.
func NewDefaultRouter(log logger.Logger) routing.Router {
	return routing.NewGorillaMux(log)
}

// New creates a new service.
func New(name, mode string, cfg config.Config, log logger.Logger, router routing.Router) *Service {
	if mode != ModeDevelopment && mode != ModeProduction {
		mode = ModeDevelopment
	}

	svc := Service{
		name:   name,
		mode:   mode,
		cfg:    cfg,
		log:    log,
		router: router,
	}

	svc.server = &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("address"),
		ReadTimeout:  time.Duration(cfg.GetInt("readTimeout")) * time.Second,
		WriteTimeout: time.Duration(cfg.GetInt("writeTimeout")) * time.Second,
	}

	return &svc
}

// NewDefault creates a new Service instance using default configuration.
func NewDefault(name string) *Service {
	cfg := NewDefaultConfig(name)
	mode := cfg.GetString("mode")
	log := NewDefaultLogger(mode)
	router := NewDefaultRouter(log)
	svc := New(name, mode, cfg, log, router)

	router.AddMiddleware(svc.ServerSignatureMiddleware)
	if mode == ModeDevelopment {
		router.AddMiddleware(svc.RequestLogDebugMiddleware)
	}

	return svc
}

// NewTesting creates a new service instance suitable to be a part of unit tests.
func NewTesting(name string) *Service {
	log := logger.NewMemory(logger.LDebug)
	svc := New(
		name,
		ModeDevelopment,
		config.NewMemory(),
		log,
		NewDefaultRouter(log),
	)

	return svc
}
