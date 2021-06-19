// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"log"

	"ampho.xyz/ampho/service"
)

// This example shows how to instantiate and run a service using default configuration.
func Example_newDefault() {
	// Create a service instance
	svc, err := service.NewDefault("hello")
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
	// Create a service instance
	svc, err := service.NewTesting("hello")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Service setup code here
	// ...

	// Start the service
	go svc.Start()

	// Base testing code here
	// ...

	// Stop the service
	svc.Stop()
}
