// Package config loads DevOS kernel configuration from environment variables
// (with an optional prefix) and optional .env-style files.
//
// It is intentionally dependency-free (standard library only) so the kernel
// can bootstrap without external configuration libraries. The shared
// packages/go/config library (a later Sprint 0 deliverable) may extend this
// for tenant services; this package covers kernel bootstrap needs.
//
// See governance/05-coding-standards.md and SDD §11.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds kernel-level configuration.
type Config struct {
	Environment     string
	ServiceName     string
	LogLevel        string
	LogFormat       string
	OTelEnabled     bool
	OTelServiceName string
	OTelEndpoint    string
	OTelProtocol    string
	NATSURL         string
	Extra           map[string]string
}

// Option configures a Loader.
type Option func(*loader)

type loader struct {
	prefix string
	files  []string
}

// WithEnvPrefix sets the environment variable prefix (default "DEVOS_").
func WithEnvPrefix(prefix string) Option {
	return func(l *loader) { l.prefix = prefix }
}

// WithFile adds an .env-style file to read before environment variables.
// Files are read in order; environment variables take precedence over files.
func WithFile(path string) Option {
	return func(l *loader) { l.files = append(l.files, path) }
}

// prefix is the default environment variable prefix.
const prefix = "DEVOS_"

// Load reads configuration using the provided options and applies defaults
// and validation.
func Load(opts ...Option) (*Config, error) {
	l := &loader{prefix: prefix}
	for _, o := range opts {
		o(l)
	}

	values := l.readFiles()

	// Environment overrides file values for prefixed keys.
	for _, kv := range os.Environ() {
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			key := kv[:idx]
			if strings.HasPrefix(key, l.prefix) {
				values[key] = kv[idx+1:]
			}
		}
	}

	cfg := &Config{
		Environment:     "development",
		ServiceName:     "devos",
		LogLevel:        "info",
		LogFormat:       "json",
		OTelEnabled:     true,
		OTelServiceName: "devos",
		OTelEndpoint:    "http://localhost:4317",
		OTelProtocol:    "grpc",
		NATSURL:         "nats://localhost:4222",
		Extra:           map[string]string{},
	}

	if v, ok := values[l.prefix+"ENV"]; ok && v != "" {
		cfg.Environment = v
	}
	if v, ok := values[l.prefix+"SERVICE_NAME"]; ok && v != "" {
		cfg.ServiceName = v
	}
	if v, ok := values[l.prefix+"LOG_LEVEL"]; ok && v != "" {
		cfg.LogLevel = strings.ToLower(v)
	}
	if v, ok := values[l.prefix+"LOG_FORMAT"]; ok && v != "" {
		cfg.LogFormat = strings.ToLower(v)
	}
	if v, ok := values[l.prefix+"OTEL_ENABLED"]; ok && v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("config: invalid %sOTEL_ENABLED: %w", l.prefix, err)
		}
		cfg.OTelEnabled = b
	}
	if v, ok := values[l.prefix+"OTEL_SERVICE_NAME"]; ok && v != "" {
		cfg.OTelServiceName = v
	}
	if v, ok := values[l.prefix+"OTEL_ENDPOINT"]; ok && v != "" {
		cfg.OTelEndpoint = v
	}
	if v, ok := values[l.prefix+"OTEL_PROTOCOL"]; ok && v != "" {
		cfg.OTelProtocol = strings.ToLower(v)
	}
	if v, ok := values[l.prefix+"NATS_URL"]; ok && v != "" {
		cfg.NATSURL = v
	}

	for k, v := range values {
		if !strings.HasPrefix(k, l.prefix) {
			continue
		}
		switch k {
		case l.prefix + "ENV", l.prefix + "SERVICE_NAME", l.prefix + "LOG_LEVEL",
			l.prefix + "LOG_FORMAT", l.prefix + "OTEL_ENABLED", l.prefix + "OTEL_SERVICE_NAME",
			l.prefix + "OTEL_ENDPOINT", l.prefix + "OTEL_PROTOCOL", l.prefix + "NATS_URL":
			// already applied above
		default:
			cfg.Extra[strings.TrimPrefix(k, l.prefix)] = v
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate enforces allowed value ranges.
func (c *Config) Validate() error {
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("config: invalid LOG_LEVEL %q (want debug|info|warn|error)", c.LogLevel)
	}
	switch c.LogFormat {
	case "json", "text":
	default:
		return fmt.Errorf("config: invalid LOG_FORMAT %q (want json|text)", c.LogFormat)
	}
	if c.Environment == "" {
		return errors.New("config: ENV must not be empty")
	}
	return nil
}

// readFiles parses .env-style files into a key/value map. Lines beginning
// with '#' and blank lines are ignored.
func (l *loader) readFiles() map[string]string {
	out := map[string]string{}
	for _, path := range l.files {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		func() {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := strings.TrimSpace(sc.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				if idx := strings.IndexByte(line, '='); idx >= 0 {
					key := strings.TrimSpace(line[:idx])
					val := strings.TrimSpace(line[idx+1:])
					out[key] = val
				}
			}
		}()
	}
	return out
}
