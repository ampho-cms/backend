// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"log"

	"ampho.xyz/core/config"
	"ampho.xyz/core/service"
)

type AwesomeService struct {
	service.Service
}

// This example shows how to extend a base service.
func Example_extend() {
	const svName = "hello"

	// Config
	cfg, err := config.New(svName, "yaml", config.DefaultSearchPaths()...)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	// Create a base service instance
	base, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Embed base service
	svc := AwesomeService{base}

	// Base setup code here
	// ...

	// Run the service until SIGINT
	svc.Run()
}
