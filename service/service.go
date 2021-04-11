// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

import (
	"ampho/routing"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"ampho/config"
	"ampho/logger"
)

// Service is the service.
type Service struct {
	name   string
	mode   string
	cfg    config.Config
	log    logger.Logger
	router routing.Router
	server *http.Server
}

// Name returns service name.
func (s *Service) Name() string {
	return s.name
}

// Mode returns mode.
func (s *Service) Mode() string {
	return s.mode
}

// Cfg returns configurator.
func (s *Service) Cfg() config.Config {
	return s.cfg
}

// Log returns logger.
func (s *Service) Log() logger.Logger {
	return s.log
}

// Router returns service router.
func (s *Service) Router() routing.Router {
	return s.router
}

// Server returns service server.
func (s *Service) Server() *http.Server {
	return s.server
}

// Start starts the service.
func (s *Service) Start() {
	defer func(logger logger.Logger) {
		if err := logger.Sync(); err != nil {
			fmt.Printf("error while syncing logger: %s", err)
		}
	}(s.log)

	// Run server without blocking to let further code process OS signals properly
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	s.log.InfoF("service '%s' started at http://%s", s.name, s.server.Addr)

	// Graceful shutdown when quit via SIGINT (Ctrl+C),
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Wait while all connections will be finished
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DftShutdownTimeout)*time.Second)
	defer cancel()
	s.log.Info("shutting down")
	_ = s.server.Shutdown(ctx)

	os.Exit(0)
}

func newDefaultLogger(mode string) logger.Logger {
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

	zpCfg.DisableCaller = true // disable caller logging since it doesn't work correctly
	zp, err = zpCfg.Build()
	if err != nil {
		panic(fmt.Errorf("Fatal error while initializing a logger: %s\n", err))
	}
	lg := logger.NewZap(zp.Sugar())
	lg.InfoF("logging initialized in %s mode", mode)

	return lg
}

func newDefaultConfig(name string) config.Config {
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

	// Load configuration
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("Fatal error while loading config file: %s \n", err))
		}
	}

	return cfg
}

func newDefaultRouter(log logger.Logger) routing.Router {
	return routing.NewGorillaMux(log)
}

// New creates a new service.
func New(name, mode string, cfg config.Config, log logger.Logger, router routing.Router) *Service {
	rand.Seed(time.Now().UnixNano())

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

	router.AddMiddleware(svc.loggingMiddleware)
	if mode == ModeDevelopment {
		router.AddMiddleware(svc.serverSignatureMiddleware)
	}
	log.Debug("router initialized")

	/// HTTP server
	svc.server = &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("address"),
		ReadTimeout:  time.Duration(cfg.GetInt("readTimeout")) * time.Second,
		WriteTimeout: time.Duration(cfg.GetInt("writeTimeout")) * time.Second,
	}
	log.Debug("server initialized")

	return &svc
}

// NewDefault creates a new service instance using default configuration.
func NewDefault(name string) *Service {
	cfg := newDefaultConfig(name)
	mode := cfg.GetString("mode")
	log := newDefaultLogger(mode)
	router := newDefaultRouter(log)

	return New(name, mode, cfg, log, router)
}
