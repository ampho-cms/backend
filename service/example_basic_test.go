// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"log"

	"ampho.xyz/ampho/routing"
	"ampho.xyz/ampho/service"
)

// helloHandler is a request handler.
func helloHandler(req *routing.Request, resp *routing.Response) {
	_, _ = resp.WriteJSON(map[string]interface{}{
		"name": req.Var("name"),
	})
}

// This example shows how to instantiate and run a simple service using default configuration.
func Example_basic() {
	// Create service instance
	svc, err := service.NewDefault("hello")
	if err != nil {
		log.Fatal(err)
	}

	// Attach a request handler
	svc.Router().Handle("/hello/{name}", helloHandler)

	// Run the service
	svc.Run()
}
