// Package lifecycle provides ordered component startup, graceful shutdown,
// and health aggregation for the DevOS kernel.
//
// Components are registered in order; Start runs Init then Start for each in
// registration order, rolling back (Stop in reverse) on the first failure.
// Stop runs in reverse registration order. Run blocks until the context is
// cancelled and then stops. This is infrastructure only — no business logic.
//
// See SDD §11 (Kernel lifecycle state machine) and Engineering Standards §6.
package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Phase is the lifecycle phase of the manager.
type Phase int

const (
	// PhaseIdle means no Start has been attempted.
	PhaseIdle Phase = iota
	// PhaseStarting means Start is in progress.
	PhaseStarting
	// PhaseRunning means all components started successfully.
	PhaseRunning
	// PhaseStopping means Stop is in progress.
	PhaseStopping
	// PhaseStopped means all components stopped.
	PhaseStopped
	// PhaseFailed means a Start failed and rollback occurred.
	PhaseFailed
)

// String renders the phase name.
func (p Phase) String() string {
	switch p {
	case PhaseIdle:
		return "idle"
	case PhaseStarting:
		return "starting"
	case PhaseRunning:
		return "running"
	case PhaseStopping:
		return "stopping"
	case PhaseStopped:
		return "stopped"
	case PhaseFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// HealthStatus is the reported status of a component.
type HealthStatus string

const (
	// StatusUp means the component is healthy.
	StatusUp HealthStatus = "UP"
	// StatusDown means the component is not healthy.
	StatusDown HealthStatus = "DOWN"
	// StatusDegraded means the component is functional but impaired.
	StatusDegraded HealthStatus = "DEGRADED"
	// StatusUnknown means health is not yet determined.
	StatusUnknown HealthStatus = "UNKNOWN"
)

// Health is a point-in-time health report for a component.
type Health struct {
	Status  HealthStatus
	Message string
	Since   time.Time
}

// Component is anything the lifecycle manager can start, stop, and probe.
type Component interface {
	Name() string
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health() Health
}

// ComponentHealth pairs a component name with its current health.
type ComponentHealth struct {
	Name   string
	Health Health
}

// Manager owns an ordered set of components and their lifecycle.
type Manager struct {
	mu         sync.RWMutex
	components []Component
	phase      Phase
}

// NewManager creates an empty manager.
func NewManager() *Manager {
	return &Manager{phase: PhaseIdle}
}

// Register appends a component in startup order.
func (m *Manager) Register(c Component) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.components = append(m.components, c)
}

// Components returns the registered components in registration order.
func (m *Manager) Components() []Component {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Component, len(m.components))
	copy(out, m.components)
	return out
}

// Phase returns the current lifecycle phase.
func (m *Manager) Phase() Phase {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.phase
}

// Init initializes all components in registration order without starting them.
func (m *Manager) Init(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.phase = PhaseStarting
	for i := 0; i < len(m.components); i++ {
		if err := m.components[i].Init(ctx); err != nil {
			m.phase = PhaseFailed
			m.rollbackLocked(ctx, i)
			return fmt.Errorf("lifecycle: init %s: %w", m.components[i].Name(), err)
		}
	}
	return nil
}

// Start initializes (if needed) and starts all components in order. On the
// first failure it stops the already-started prefix in reverse and returns.
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.phase = PhaseStarting
	for i := 0; i < len(m.components); i++ {
		c := m.components[i]
		if err := c.Init(ctx); err != nil {
			m.phase = PhaseFailed
			m.rollbackLocked(ctx, i)
			return fmt.Errorf("lifecycle: init %s: %w", c.Name(), err)
		}
		if err := c.Start(ctx); err != nil {
			m.phase = PhaseFailed
			m.rollbackLocked(ctx, i)
			return fmt.Errorf("lifecycle: start %s: %w", c.Name(), err)
		}
	}
	m.phase = PhaseRunning
	return nil
}

// Stop stops all components in reverse registration order.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stopAllLocked(ctx)
}

// Health returns the current health of every component.
func (m *Manager) Health() []ComponentHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ComponentHealth, 0, len(m.components))
	for _, c := range m.components {
		out = append(out, ComponentHealth{Name: c.Name(), Health: c.Health()})
	}
	return out
}

// Run starts the manager and blocks until ctx is cancelled, then stops all
// components with a bounded graceful-shutdown timeout.
func (m *Manager) Run(ctx context.Context) error {
	if err := m.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return m.stopAllLocked(stopCtx)
}

// rollbackLocked stops components [0, upto) in reverse order.
func (m *Manager) rollbackLocked(ctx context.Context, upto int) {
	for j := upto - 1; j >= 0; j-- {
		_ = m.components[j].Stop(ctx)
	}
}

// stopAllLocked stops all components in reverse registration order,
// returning the first error encountered.
func (m *Manager) stopAllLocked(ctx context.Context) error {
	m.phase = PhaseStopping
	var firstErr error
	for i := len(m.components) - 1; i >= 0; i-- {
		if err := m.components[i].Stop(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	m.phase = PhaseStopped
	return firstErr
}
