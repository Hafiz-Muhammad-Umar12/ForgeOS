// Package kernel is the DevOS kernel composition root. It loads configuration,
// builds the dependency-injection container, registers core services, and
// constructs the lifecycle manager.
//
// This is the skeleton: it wires DI + lifecycle + config and exposes the
// runtime facade. Bus/registry/scheduler client implementations (NATS, etc.)
// and tenant services are attached in later milestones (SDD §11, build order
// Step 1+). No business logic, network I/O, or authentication is performed
// here beyond configuration loading.
package kernel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/config"
	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/di"
	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/lifecycle"
)

// Kernel is the composed runtime. It owns the configuration, the
// dependency-injection container, and the lifecycle manager.
type Kernel struct {
	Config    *config.Config
	Container *di.Container
	Lifecycle *lifecycle.Manager
}

// Option configures kernel construction.
type Option func(*options)

type options struct {
	configOpts []config.Option
}

// WithConfigOpts forwards configuration options to config.Load.
func WithConfigOpts(opts ...config.Option) Option {
	return func(o *options) { o.configOpts = append(o.configOpts, opts...) }
}

// New loads configuration and wires the core services into a DI container and
// lifecycle manager. It returns a ready-to-run kernel.
func New(opts ...Option) (*Kernel, error) {
	o := options{}
	for _, fn := range opts {
		fn(&o)
	}

	cfg, err := config.Load(o.configOpts...)
	if err != nil {
		return nil, fmt.Errorf("kernel: load config: %w", err)
	}

	container := di.New()
	container.MustRegister(reflect.TypeOf(&config.Config{}), func(*di.Container) (any, error) {
		return cfg, nil
	}, di.Singleton)
	container.MustRegister(reflect.TypeOf(&di.Container{}), func(*di.Container) (any, error) {
		return container, nil
	}, di.Singleton)

	lm := lifecycle.NewManager()

	return &Kernel{
		Config:    cfg,
		Container: container,
		Lifecycle: lm,
	}, nil
}

// Run starts the kernel and blocks until ctx is cancelled, then stops all
// registered components gracefully.
func (k *Kernel) Run(ctx context.Context) error {
	return k.Lifecycle.Run(ctx)
}
