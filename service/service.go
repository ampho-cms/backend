// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"ampho.xyz/ampho/config"
	"ampho.xyz/ampho/logging"
	"ampho.xyz/ampho/routing"
	"ampho.xyz/ampho/security"
)

// Service is the service interface.
type Service interface {
	Name() string
	Signature() string
	Router() routing.Router
	Server() *http.Server
	Start()
	Stop()
	OnStop(fn func(svc Service))
	Run()
}

// Base is the base service structure.
type Base struct {
	name   string
	router routing.Router
	server *http.Server
	onStop []func(svc Service)
}

// Name returns service name.
func (s *Base) Name() string {
	return s.name
}

// Signature returns service signature, usually used as a 'Server' HTTP header.
func (s *Base) Signature() string {
	return s.name
}

// Router returns service router.
func (s *Base) Router() routing.Router {
	return s.router
}

// Server returns service HTTP server.
func (s *Base) Server() *http.Server {
	return s.server
}

// Start starts the service. Assumed to be run as a goroutine.
func (s *Base) Start() {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		logging.InfoF("%s has been stopped", s.name)
	} else {
		logging.Error(err)
	}
}

// OnStop schedules a function to be called at service stop phase.
func (s *Base) OnStop(fn func(svc Service)) {
	s.onStop = append(s.onStop, fn)
}

// Stop stops the service.
func (s *Base) Stop() {
	logging.InfoF("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), config.GetDuration("service.shutdownTimeout"))
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logging.Error(err.Error())
	}

	for _, fn := range s.onStop {
		fn(s)
	}
}

// Run starts the service and waits for SIGINT. Usually it should be a last call in the `main()`.
func (s *Base) Run() {
	go s.Start()
	logging.InfoF("%s started at http://%s", s.name, s.server.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()
}

// New creates a new base service instance.
func New(name string, router routing.Router, srv *http.Server) *Base {
	return &Base{name, router, srv, nil}
}

// configDefaults returns service configuration defaults.
func configDefaults() map[string]interface{} {
	return map[string]interface{}{
		"logging.level":           DftLoggingLevel,
		"logging.file.size":       DftLoggingFileSize,
		"logging.file.age":        DftLoggingFileAge,
		"logging.file.backups":    DftLoggingFileBackups,
		"logging.file.stdout":     DftLoggingFileStdout,
		"security.signingMethod":  DftSecuritySigningMethod,
		"service.address":         DftAddress,
		"service.readTimeout":     DftReadTimeout,
		"service.writeTimeout":    DftWriteTimeout,
		"service.shutdownTimeout": DftShutdownTimeout,
	}
}

func initConfig(name string) error {
	// Configuration engine
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	cfgSearchPaths := []string{"$HOME/." + name, execDir, "."}
	cfg, err := config.NewDefault(name, "yaml", cfgSearchPaths, configDefaults())
	if err != nil {
		return err
	}

	config.SetConfig(cfg)

	return nil
}

// initLogging initializes a logger from a config.
func initLogging(name string) (logging.Logger, error) {
	var (
		err    error
		logger logging.Logger
	)

	level := logging.GetLevelByName(config.GetString("logging.level"))
	path := config.GetString("logging.file.path")

	if path != "" {
		path = strings.ReplaceAll(path, "$SERVICE_NAME", name)
		size := config.GetInt("logging.file.size")
		age := config.GetInt("logging.file.age")
		backups := config.GetInt("logging.file.backups")
		stdout := config.GetBool("logging.file.stdout")
		if logger, err = logging.NewFileLogger(path, size, age, backups, level, stdout); err != nil {
			return nil, err
		}
	} else {
		if logger, err = logging.NewConsoleLogger(level); err != nil {
			return nil, err
		}
	}

	logging.SetLogger(logger)

	return logger, nil
}

// initSecurity initializes security package from config.
func initSecurity() error {
	var err error

	// Signing method
	if err = security.SetSigningMethod(config.GetString("security.signingMethod")); err != nil {
		return fmt.Errorf("failed to set signing method: %v", err)
	}

	// HMAC key
	hmacKey := config.GetString("security.hmac.key")
	if hmacKey != "" {
		security.SetHMACKey([]byte(hmacKey))
	}

	// RSA keys
	rsaPrvKey := config.GetString("security.rsa.privateKey")
	if rsaPrvKey != "" {
		if err = security.SetRSAPrivateKey([]byte(rsaPrvKey)); err != nil {
			return fmt.Errorf("failed to load private RSA key: %v", err)
		}
	}
	rsaPubKey := config.GetString("security.rsa.publicKey")
	if rsaPubKey != "" {
		if err = security.SetRSAPublicKey([]byte(rsaPubKey)); err != nil {
			return fmt.Errorf("failed to load public RSA key: %v", err)
		}
	}

	// ECDSA keys
	ecdsaPrvKey := config.GetString("security.ecdsa.privateKey")
	if ecdsaPrvKey != "" {
		if err = security.SetECDSAPrivateKey([]byte(ecdsaPrvKey)); err != nil {
			return fmt.Errorf("failed to load private ECDSA key: %v", err)
		}
	}
	ecdsaPubKey := config.GetString("security.ecdsa.publicKey")
	if ecdsaPubKey != "" {
		if err = security.SetECDSAPublicKey([]byte(ecdsaPubKey)); err != nil {
			return fmt.Errorf("failed to load public ECDSA key: %v", err)
		}
	}

	return nil
}

// NewDefault creates a new base service instance using default configuration.
func NewDefault(name string) (*Base, error) {
	var err error

	// Config
	if err = initConfig(name); err != nil {
		return nil, fmt.Errorf("failed to init config: %v", err)
	}

	// Logging
	logger, err := initLogging(name)
	if err != nil {
		return nil, fmt.Errorf("failed to init logging: %v", err)
	}

	// Security
	if err = initSecurity(); err != nil {
		return nil, fmt.Errorf("failed to init security: %v", err)
	}

	// Router
	router := routing.NewDefault()

	// Server
	srv := &http.Server{
		Handler:      router,
		Addr:         config.GetString("service.address"),
		ReadTimeout:  config.GetDuration("service.readTimeout"),
		WriteTimeout: config.GetDuration("service.writeTimeout"),
	}

	// Service
	svc := New(name, router, srv)

	// Flush logs on service stop
	svc.OnStop(func(svc Service) {
		logging.DebugF("flushing logs...")
		_ = logger.Sync()
	})

	// Middlewares
	router.AddMiddleware(svc.ServerSignatureMiddleware)
	if logger.Level() == logging.LDebug {
		router.AddMiddleware(svc.RequestLogDebugMiddleware)
	}

	return svc, nil
}

// NewTesting creates a new base service instance suitable for unit tests.
func NewTesting(name string) (*Base, error) {
	// Config
	cfg, err := config.NewDefault(name, "", nil, configDefaults())
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %v", err)
	}
	config.SetConfig(cfg)

	// Logging
	logger, err := logging.NewConsoleLogger(logging.LDebug)
	if err != nil {
		return nil, fmt.Errorf("failed to init logging: %v", err)
	}
	logging.SetLogger(logger)

	// Security
	if err = initSecurity(); err != nil {
		return nil, fmt.Errorf("failed to init security: %v", err)
	}

	// Router
	router := routing.NewDefault()

	// Server
	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("service.address"),
		ReadTimeout:  cfg.GetDuration("service.readTimeout"),
		WriteTimeout: cfg.GetDuration("service.writeTimeout"),
	}

	return New(name, router, srv), nil
}
