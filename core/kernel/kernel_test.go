package kernel

import (
	"context"
	"testing"
)

func TestNewDefaults(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if k == nil || k.Config == nil || k.Container == nil || k.Lifecycle == nil {
		t.Fatal("nil fields")
	}
	if k.Config.Environment != "development" {
		t.Errorf("env=%s", k.Config.Environment)
	}
}

func TestNewConfigOverride(t *testing.T) {
	t.Setenv("DEVOS_ENV", "staging")
	k, err := New(WithConfigOpts())
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if k.Config.Environment != "staging" {
		t.Errorf("env=%s", k.Config.Environment)
	}
}

func TestRunStopsOnCancelledContext(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := k.Run(ctx); err != nil {
		t.Fatalf("run: %v", err)
	}
}
