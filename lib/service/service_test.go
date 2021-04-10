// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service_test

import (
	"testing"

	"ampho/service"
)

// TestNewDefault tests default service creation
func TestNewDefault(t *testing.T) {
	name := "HelloWorld"
	app := service.NewDefault("HelloWorld")
	got := app.Name()
	if got != name {
		t.Errorf("Name() == %q, want %q", got, name)
	}
}
