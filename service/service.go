// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
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
	router *mux.Router
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
func (s *Service) Router() *mux.Router {
	return s.router
}

// Server returns service server.
func (s *Service) Server() *http.Server {
	return s.server
}

// AddHandler registers a request handler.
func (s *Service) AddHandler(path string, handler RequestHandler) *mux.Route {
	route := s.router.HandleFunc(path, func(resp http.ResponseWriter, req *http.Request) {
		handler(s, &Request{request: req}, &Response{writer: resp})
	})

	s.log.DebugF("route registered: %s", path)

	return route
}

// AddMiddleware registers a middleware.
func (s *Service) AddMiddleware(handler RequestHandler) {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			handler(s, &Request{request: req}, &Response{writer: resp})
			next.ServeHTTP(resp, req)
		})
	})

	s.log.DebugF("middleware registered: %v", handler)
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

// NewDefault creates a new service instance using default configuration.
func NewDefault(name string) *Service {
	var err error

	rand.Seed(time.Now().UnixNano())

	s := Service{
		name: name,
		mode: ModeDevelopment,
	}

	// Configuration engine
	vp := viper.New()
	vp.SetConfigName("." + s.name + ".yaml")
	vp.AddConfigPath(".")
	s.cfg = config.NewViper(vp)

	// Configuration defaults
	s.cfg.SetDefault("mode", ModeDevelopment)
	s.cfg.SetDefault("address", DftNetAddr)
	s.cfg.SetDefault("readTimeout", DftNetReadTimeout)
	s.cfg.SetDefault("writeTimeout", DftNetWriteTimeout)

	// Load configuration
	if err = vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("Fatal error while loading config file: %s \n", err))
		}
	}

	// Update mode from config
	mode := s.cfg.GetString("mode")
	if mode != ModeDevelopment && mode != ModeProduction {
		panic(fmt.Errorf("unknown mode: %s", mode))
	}
	s.mode = mode

	// Logging
	var (
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
	s.log = logger.NewZap(zp.Sugar())
	s.log.InfoF("logging initialized in %s mode", mode)

	// Router
	s.router = mux.NewRouter()
	s.AddMiddleware(loggingMiddleware)
	if mode == ModeDevelopment {
		s.AddMiddleware(serverSignatureMiddleware)
	}
	s.log.Debug("router initialized")

	/// HTTP server
	s.server = &http.Server{
		Handler:      s.router,
		Addr:         s.cfg.GetString("address"),
		ReadTimeout:  time.Duration(s.cfg.GetInt("readTimeout")) * time.Second,
		WriteTimeout: time.Duration(s.cfg.GetInt("writeTimeout")) * time.Second,
	}
	s.log.Debug("server initialized")

	return &s
}
