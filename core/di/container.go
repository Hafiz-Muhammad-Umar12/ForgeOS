// Package di provides a small, reflection-based dependency-injection
// container used as the composition root for the DevOS kernel.
//
// It supports registering factories under a contract type, singleton and
// transient lifetimes, resolution by type, and circular-dependency detection.
// It is intentionally minimal: it wires the kernel and tenant services at
// startup, not as a general-purpose DI framework.
//
// See governance/04-engineering-standards.md §3 (Dependency Injection) and
// SDD §11 (Kernel microkernel).
package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrNotRegistered is returned when resolving an unregistered contract.
var ErrNotRegistered = errors.New("di: contract not registered")

// ErrAlreadyRegistered is returned when registering a duplicate contract.
var ErrAlreadyRegistered = errors.New("di: contract already registered")

// ErrCircularDependency is returned when a dependency cycle is detected.
var ErrCircularDependency = errors.New("di: circular dependency detected")

// ErrFactoryPanic is returned when a factory function panics.
var ErrFactoryPanic = errors.New("di: factory panicked")

// Lifetime controls how instances are cached.
type Lifetime int

const (
	// Singleton creates the instance once and reuses it for every resolve.
	Singleton Lifetime = iota
	// Transient creates a new instance on every resolve.
	Transient
)

// Factory builds an instance, optionally resolving its own dependencies
// from the container.
type Factory func(c *Container) (any, error)

type registration struct {
	factory  Factory
	lifetime Lifetime
	instance any
	built    bool
	building bool
}

// Container maps contract types to factories.
type Container struct {
	mu   sync.RWMutex
	regs map[reflect.Type]*registration
}

// New creates an empty container.
func New() *Container {
	return &Container{regs: make(map[reflect.Type]*registration)}
}

// Register binds factory to contract. contract must be a non-nil
// reflect.Type, typically obtained via reflect.TypeOf((*T)(nil)).Elem() for an
// interface T, or reflect.TypeOf(T{}) for a struct T.
func (c *Container) Register(contract reflect.Type, factory Factory, lifetime Lifetime) error {
	if contract == nil {
		return errors.New("di: contract type is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.regs[contract]; exists {
		return fmt.Errorf("%w: %s", ErrAlreadyRegistered, contract)
	}
	c.regs[contract] = &registration{factory: factory, lifetime: lifetime}
	return nil
}

// MustRegister is Register that panics on error. Use it in composition roots
// where a registration failure is fatal.
func (c *Container) MustRegister(contract reflect.Type, factory Factory, lifetime Lifetime) {
	if err := c.Register(contract, factory, lifetime); err != nil {
		panic(err)
	}
}

// Has reports whether contract is registered.
func (c *Container) Has(contract reflect.Type) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.regs[contract]
	return ok
}

// Resolve returns the instance bound to contract.
func (c *Container) Resolve(contract reflect.Type) (any, error) {
	return c.resolve(contract)
}

// MustResolve is Resolve that panics on error.
func (c *Container) MustResolve(contract reflect.Type) any {
	v, err := c.Resolve(contract)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Container) resolve(contract reflect.Type) (any, error) {
	c.mu.RLock()
	reg, ok := c.regs[contract]
	c.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotRegistered, contract)
	}
	if reg.lifetime == Singleton && reg.built {
		return reg.instance, nil
	}

	c.mu.Lock()
	if reg.building {
		c.mu.Unlock()
		return nil, fmt.Errorf("%w: %s", ErrCircularDependency, contract)
	}
	reg.building = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		reg.building = false
		c.mu.Unlock()
	}()

	var (
		instance any
		err      error
	)
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%w: %v", ErrFactoryPanic, r)
			}
		}()
		instance, err = reg.factory(c)
	}()
	if err != nil {
		return nil, err
	}

	got := reflect.TypeOf(instance)
	if got == nil || !got.AssignableTo(contract) {
		return nil, fmt.Errorf("di: factory returned %s, not assignable to %s", got, contract)
	}

	if reg.lifetime == Singleton {
		c.mu.Lock()
		reg.instance = instance
		reg.built = true
		c.mu.Unlock()
	}
	return instance, nil
}
