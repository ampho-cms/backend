// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service_test provides tests of base service structures and functions.
package service

import (
	"testing"

	"ampho.xyz/ampho/util"
)

// TestNewDefault tests default service creation.
func TestNewDefault(t *testing.T) {
	name := util.RandAscii(8)
	svc, err := NewTesting(name)
	if err != nil {
		t.Errorf("%v", err)
	}

	if got := svc.Name(); got != name {
		t.Errorf("Name() == %q, want %q", got, name)
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
