// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"log"

	"ampho.xyz/ampho/service"
)

type AwesomeService struct {
	service.Service
}

// This example shows how to extend a base service.
func Example_extend() {
	// Create a base service instance
	base, err := service.NewDefault("hello")
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
