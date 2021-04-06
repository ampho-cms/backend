package service_test

import (
	"testing"

	"ampho/service"
)

func TestNew(t *testing.T) {
	name := "HelloWorld"
	app := service.New("HelloWorld")
	got := app.Name()
	if got != name {
		t.Errorf("Name() == %q, want %q", got, name)
	}
}
