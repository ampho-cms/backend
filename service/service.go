// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"

	"ampho.xyz/config"
)

// Service is the service interface.
type Service interface {
	// Config returns configuration.
	Config() config.Config

	// Router returns router.
	Router() *mux.Router

	// Server returns HTTP server.
	Server() *http.Server

	// BeforeStart schedules a function to be called before service start.
	BeforeStart(fn func(svc Service))

	// AfterStop schedules a function to be called after service stop.
	AfterStop(fn func(svc Service))

	// Start starts the service. Assumed to be called in a goroutine.
	Start()

	// Stop stops the service.
	Stop()

	// Run starts the service and blocks until SIGINT received. Usually it should be a last call in the `main()`.
	Run()
}

// Base is the base service structure.
type Base struct {
	config      config.Config
	server      *http.Server
	beforeStart []func(svc Service)
	afterStop   []func(svc Service)
}

// Config returns configuration.
func (s *Base) Config() config.Config {
	return s.config
}

// Router returns router.
func (s *Base) Router() *mux.Router {
	return s.server.Handler.(*mux.Router)
}

// Server returns HTTP server.
func (s *Base) Server() *http.Server {
	return s.server
}

// BeforeStart schedules a function to be called before service start.
func (s *Base) BeforeStart(fn func(svc Service)) {
	s.beforeStart = append(s.beforeStart, fn)
}

// AfterStop schedules a function to be called after service stop.
func (s *Base) AfterStop(fn func(svc Service)) {
	s.afterStop = append(s.afterStop, fn)
}

// Start starts the service. Assumed to be called in a goroutine.
func (s *Base) Start() {
	for _, fn := range s.beforeStart {
		fn(s)
	}

	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Stop stops the service.
func (s *Base) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.Config().GetDuration("service.shutdownTimeout"))
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Panic(err.Error())
	}

	for _, fn := range s.afterStop {
		fn(s)
	}
}

// Run starts the service and blocks until SIGINT received. Usually it should be a last call in the `main()`.
func (s *Base) Run() {
	go s.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()
}

// New creates a new base service instance using default configuration.
func New(cfg config.Config) (*Base, error) {
	cfg.SetDefault("service.address", DftAddress)
	cfg.SetDefault("service.readTimeout", DftReadTimeout)
	cfg.SetDefault("service.writeTimeout", DftWriteTimeout)
	cfg.SetDefault("service.shutdownTimeout", DftShutdownTimeout)

	// Server
	srv := &http.Server{
		Handler:      mux.NewRouter(),
		Addr:         cfg.GetString("service.address"),
		ReadTimeout:  cfg.GetDuration("service.readTimeout"),
		WriteTimeout: cfg.GetDuration("service.writeTimeout"),
	}

	// Service
	svc := &Base{cfg, srv, nil, nil}

	// Middlewares
	svc.Router().Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	return svc, nil
}
