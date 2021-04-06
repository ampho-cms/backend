// Package service provides base service structures and functions
package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Service is the service.
type Service struct {
	name   string
	config *viper.Viper
	logger *zap.SugaredLogger
	router *mux.Router
	server *http.Server
}


// Name returns service name.
func (s *Service) Name() string {
	return s.name
}

// Config returns service configurator.
func (s *Service) Config() *viper.Viper {
	return s.config
}

// Logger returns service logger.
func (s *Service) Logger() *zap.SugaredLogger {
	return s.logger
}

// Router returns service router.
func (s *Service) Router() *mux.Router {
	return s.router
}

// Server returns service server.
func (s *Service) Server() *http.Server {
	return s.server
}

// Handle registers a new route for the URL path.
func (s *Service) Handle(path string, handler RequestHandler) *mux.Route {
	return s.router.HandleFunc(path, func(resp http.ResponseWriter, req *http.Request) {
		handler(s, &Request{request: req}, &Response{writer: resp})
	})
}

func (s *Service) AddMiddleware(handler RequestHandler) {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			handler(s, &Request{request: req}, &Response{writer: resp})
			next.ServeHTTP(resp, req)
		})
	})
}

// initConfig initializes a configuration subsystem and loads initial configuration.
func (s *Service) initConfig() {
	cfg := viper.New()

	cfg.SetDefault("address", defaultNetAddr)
	cfg.SetDefault("readTimeout", defaultReadTimeout)
	cfg.SetDefault("writeTimeout", defaultWriteTimeout)

	cfg.SetConfigName("." + s.name + ".yaml")
	cfg.AddConfigPath(".")

	err := cfg.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("Fatal error while loading config file: %s \n", err))
		}
	}

	s.config = cfg
}

// initLogger initializes service logging subsystem.
func (s *Service) initLogger() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("Fatal error while initializing a logger: %s\n", err))
	}

	s.logger = l.Sugar()
}

// initLogger initializes service routing subsystem.
func (s *Service) initRouter() {
	s.router = mux.NewRouter()

	s.AddMiddleware(loggingMiddleware)
	s.AddMiddleware(serverSignatureMiddleware)
}

// Start starts the service.
func (s *Service) Start() {
	defer s.logger.Sync()

	// Run server without blocking to ket further code process OS signals properly
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	s.logger.Infof("Service '%s' started at http://%s", s.name, s.server.Addr)

	// Graceful shutdown when quit via SIGINT (Ctrl+C),
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Wait while all connections will be finished
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()
	s.logger.Info("shutting down")
	_ = s.server.Shutdown(ctx)

	os.Exit(0)
}

// New creates a new service instance.
func New(name string) *Service {
	s := Service{
		name: name,
	}

	s.initConfig()
	s.initLogger()
	s.initRouter()

	s.server = &http.Server{
		Handler:      s.router,
		Addr:         s.config.GetString("address"),
		ReadTimeout:  time.Duration(s.config.GetInt("readTimeout")) * time.Second,
		WriteTimeout: time.Duration(s.config.GetInt("writeTimeout")) * time.Second,
	}

	return &s
}
