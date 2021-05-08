// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service_test provides tests of base service structures and functions.
package service

import (
	"testing"

	"ampho.xyz/ampho/util"
)

// TestNew tests service creation
func TestNew(t *testing.T) {
	name := util.RandAscii(8)
	svc := NewTesting(name)

	if got := svc.Name(); got != name {
		t.Errorf("Name() == %q, want %q", got, name)
	}

	if got := svc.Mode(); got != ModeDevelopment {
		t.Errorf("Mode() == %q, want %q", got, ModeDevelopment)
	}

	if got := svc.Cfg(); got == nil {
		t.Errorf("Cfg() == %v, want pointer", got)
	}

	if got := svc.Log(); got == nil {
		t.Errorf("Log() == %v, want pointer", got)
	}

	if got := svc.Router(); got == nil {
		t.Errorf("Router() == %v, want pointer", got)
	}

	if got := svc.Server(); got == nil {
		t.Errorf("Server() == %v, want pointer", got)
	}

	go svc.Start()
	svc.Stop()
}
