package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Environment != "development" {
		t.Errorf("env=%s", cfg.Environment)
	}
	if cfg.ServiceName != "devos" {
		t.Errorf("svc=%s", cfg.ServiceName)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("level=%s", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("fmt=%s", cfg.LogFormat)
	}
	if !cfg.OTelEnabled {
		t.Error("otel should default true")
	}
	if cfg.OTelEndpoint != "http://localhost:4317" {
		t.Errorf("otel ep=%s", cfg.OTelEndpoint)
	}
	if cfg.NATSURL != "nats://localhost:4222" {
		t.Errorf("nats=%s", cfg.NATSURL)
	}
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("DEVOS_ENV", "production")
	t.Setenv("DEVOS_LOG_LEVEL", "debug")
	t.Setenv("DEVOS_OTEL_ENABLED", "false")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Environment != "production" {
		t.Errorf("env=%s", cfg.Environment)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("level=%s", cfg.LogLevel)
	}
	if cfg.OTelEnabled {
		t.Error("otel should be false")
	}
}

func TestLoadPrefix(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	cfg, err := Load(WithEnvPrefix("APP_"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Environment != "staging" {
		t.Errorf("env=%s", cfg.Environment)
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "# comment\nDEVOS_ENV=test\nDEVOS_LOG_LEVEL=warn\nDEVOS_CUSTOM=value\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(WithFile(path))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Environment != "test" {
		t.Errorf("env=%s", cfg.Environment)
	}
	if cfg.LogLevel != "warn" {
		t.Errorf("level=%s", cfg.LogLevel)
	}
	if cfg.Extra["CUSTOM"] != "value" {
		t.Errorf("extra=%v", cfg.Extra)
	}
}

func TestLoadInvalidLogLevel(t *testing.T) {
	t.Setenv("DEVOS_LOG_LEVEL", "verbose")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid log level")
	}
}

func TestLoadInvalidLogFormat(t *testing.T) {
	t.Setenv("DEVOS_LOG_FORMAT", "yaml")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid log format")
	}
}
