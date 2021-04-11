// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service_test provides tests of base service structures and functions.
package service_test

import (
	"testing"

	"ampho/service"
	"ampho/util"
)

// TestNewDefault tests default service creation
func TestNewDefault(t *testing.T) {
	name := util.RandStrAscii(8)
	app := service.NewDefault(name)

	if got := app.Name(); got != name {
		t.Errorf("Name() == %q, want %q", got, name)
	}

	if got := app.Mode(); got != service.ModeDevelopment {
		t.Errorf("Mode() == %q, want %q", got, service.ModeDevelopment)
	}

	if got := app.Cfg(); got == nil {
		t.Errorf("Cfg() == %v, want pointer", got)
	}

	if got := app.Log(); got == nil {
		t.Errorf("Log() == %v, want pointer", got)
	}

	if got := app.Router(); got == nil {
		t.Errorf("Router() == %v, want pointer", got)
	}

	if got := app.Server(); got == nil {
		t.Errorf("Server() == %v, want pointer", got)
	}
}
