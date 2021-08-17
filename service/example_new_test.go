// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"log"

	"ampho.xyz/core/config"
	"ampho.xyz/core/service"
)

// This example shows how to instantiate and run a service using default configuration.
func Example_newDefault() {
	const svName = "hello"

	// Config
	cfg, err := config.New(svName, "yaml", config.DefaultSearchPaths()...)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	// Create a service instance
	svc, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Service setup code here
	// ...

	// Run the service until SIGINT
	svc.Run()
}

// This example shows how to instantiate and use a service using testing configuration.
func Example_newTesting() {
	// TODO
}
