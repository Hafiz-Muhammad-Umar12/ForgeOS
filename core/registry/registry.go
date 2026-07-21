// Package registry defines the service-registry port and supporting types
// for DevOS. It is the contract only; the NATS-KV-backed implementation is
// delivered in a later milestone (see SDD §09 and ADR-003).
//
// Tenant services depend on the Registry interface, never on a concrete
// implementation, so the backing store can be swapped without touching
// callers.
package registry

import (
	"context"
	"time"
)

// Capability expresses a feature a registered service provides.
type Capability string

// ServiceInfo describes a registered service, agent, or provider.
type ServiceInfo struct {
	ID           string
	Name         string
	Kind         string // "agent" | "provider" | "service" | "channel" | "tool"
	Capabilities []Capability
	Endpoint     string
	Metadata     map[string]string
	RegisteredAt time.Time
}

// Has reports whether info provides the given capability.
func (s ServiceInfo) Has(c Capability) bool {
	for _, c2 := range s.Capabilities {
		if c2 == c {
			return true
		}
	}
	return false
}

// CapabilitySet is a uniqueness-preserving collection of capabilities.
type CapabilitySet map[Capability]struct{}

// NewCapabilitySet builds a set from the given capabilities.
func NewCapabilitySet(caps ...Capability) CapabilitySet {
	s := make(CapabilitySet, len(caps))
	for _, c := range caps {
		s.Add(c)
	}
	return s
}

// Add inserts a capability.
func (s CapabilitySet) Add(c Capability) {
	s[c] = struct{}{}
}

// Contains reports whether c is in the set.
func (s CapabilitySet) Contains(c Capability) bool {
	_, ok := s[c]
	return ok
}

// Slice returns the capabilities as a slice (order is non-deterministic).
func (s CapabilitySet) Slice() []Capability {
	out := make([]Capability, 0, len(s))
	for c := range s {
		out = append(out, c)
	}
	return out
}

// Registry is the service-discovery and capability-index port.
type Registry interface {
	Register(ctx context.Context, info ServiceInfo) error
	Deregister(ctx context.Context, id string) error
	Discover(ctx context.Context, kind string) ([]ServiceInfo, error)
	Resolve(ctx context.Context, capability Capability) (ServiceInfo, error)
	ResolveByID(ctx context.Context, id string) (ServiceInfo, error)
	Watch(ctx context.Context, capability Capability) (<-chan ServiceInfo, error)
}
