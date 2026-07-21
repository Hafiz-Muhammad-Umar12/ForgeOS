package registry

import (
	"context"
	"testing"
)

func TestServiceInfoHas(t *testing.T) {
	info := ServiceInfo{Capabilities: []Capability{"llm", "embed"}}
	if !info.Has("llm") {
		t.Error("should have llm")
	}
	if info.Has("vector") {
		t.Error("should not have vector")
	}
}

func TestCapabilitySet(t *testing.T) {
	s := NewCapabilitySet("a", "b", "a")
	if len(s) != 2 {
		t.Fatalf("len=%d", len(s))
	}
	if !s.Contains("a") || !s.Contains("b") {
		t.Error("missing")
	}
	if s.Contains("c") {
		t.Error("extra")
	}
	if len(s.Slice()) != 2 {
		t.Fatalf("slice len=%d", len(s.Slice()))
	}
}

// memRegistry is an in-memory test double for the Registry port.
type memRegistry struct {
	store map[string]ServiceInfo
}

func newMem() *memRegistry { return &memRegistry{store: map[string]ServiceInfo{}} }

func (m *memRegistry) Register(_ context.Context, info ServiceInfo) error {
	m.store[info.ID] = info
	return nil
}
func (m *memRegistry) Deregister(_ context.Context, id string) error {
	delete(m.store, id)
	return nil
}
func (m *memRegistry) Discover(_ context.Context, kind string) ([]ServiceInfo, error) {
	var out []ServiceInfo
	for _, v := range m.store {
		if v.Kind == kind {
			out = append(out, v)
		}
	}
	return out, nil
}
func (m *memRegistry) Resolve(_ context.Context, cap Capability) (ServiceInfo, error) {
	for _, v := range m.store {
		if v.Has(cap) {
			return v, nil
		}
	}
	return ServiceInfo{}, nil
}
func (m *memRegistry) ResolveByID(_ context.Context, id string) (ServiceInfo, error) {
	return m.store[id], nil
}
func (m *memRegistry) Watch(_ context.Context, _ Capability) (<-chan ServiceInfo, error) {
	return make(chan ServiceInfo), nil
}

// compile-time check that the test double satisfies the port.
var _ Registry = (*memRegistry)(nil)

func TestRegistryFake(t *testing.T) {
	r := newMem()
	ctx := context.Background()
	if err := r.Register(ctx, ServiceInfo{ID: "1", Kind: "agent", Capabilities: []Capability{"plan"}}); err != nil {
		t.Fatal(err)
	}
	got, err := r.Resolve(ctx, "plan")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "1" {
		t.Errorf("resolved=%s", got.ID)
	}
	disc, _ := r.Discover(ctx, "agent")
	if len(disc) != 1 {
		t.Errorf("discover=%d", len(disc))
	}
}
