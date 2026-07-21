package lifecycle

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type recorder struct {
	mu  sync.Mutex
	log []string
}

func (r *recorder) add(s string) {
	r.mu.Lock()
	r.log = append(r.log, s)
	r.mu.Unlock()
}

func (r *recorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.log))
	copy(out, r.log)
	return out
}

type fakeComponent struct {
	name     string
	rec      *recorder
	initErr  error
	startErr error
	stopErr  error
}

func (f *fakeComponent) Name() string { return f.name }
func (f *fakeComponent) Init(context.Context) error {
	f.rec.add(f.name + ".init")
	return f.initErr
}
func (f *fakeComponent) Start(context.Context) error {
	f.rec.add(f.name + ".start")
	return f.startErr
}
func (f *fakeComponent) Stop(context.Context) error {
	f.rec.add(f.name + ".stop")
	return f.stopErr
}
func (f *fakeComponent) Health() Health {
	return Health{Status: StatusUp, Since: time.Now()}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestManagerStartStopOrder(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec})
	m.Register(&fakeComponent{name: "b", rec: rec})
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := m.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	want := []string{"a.init", "a.start", "b.init", "b.start", "b.stop", "a.stop"}
	if got := rec.snapshot(); !equal(got, want) {
		t.Fatalf("order:\n got=%v\nwant=%v", got, want)
	}
	if m.Phase() != PhaseStopped {
		t.Fatalf("phase=%v", m.Phase())
	}
}

func TestStartFailureRollback(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec})
	m.Register(&fakeComponent{name: "b", rec: rec, startErr: errors.New("fail")})
	err := m.Start(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	want := []string{"a.init", "a.start", "b.init", "b.start", "a.stop"}
	if got := rec.snapshot(); !equal(got, want) {
		t.Fatalf("rollback order:\n got=%v\nwant=%v", got, want)
	}
	if m.Phase() != PhaseFailed {
		t.Fatalf("phase=%v", m.Phase())
	}
}

func TestRunBlocksUntilCancel(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	if err := m.Run(ctx); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got := rec.snapshot(); !equal(got, []string{"a.init", "a.start", "a.stop"}) {
		t.Fatalf("run order: %v", got)
	}
}

func TestHealth(t *testing.T) {
	m := NewManager()
	m.Register(&fakeComponent{name: "a"})
	m.Register(&fakeComponent{name: "b"})
	h := m.Health()
	if len(h) != 2 {
		t.Fatalf("health len=%d", len(h))
	}
	if h[0].Name != "a" || h[1].Name != "b" {
		t.Fatalf("health names: %v", h)
	}
}

func TestComponentsOrder(t *testing.T) {
	m := NewManager()
	a := &fakeComponent{name: "a"}
	b := &fakeComponent{name: "b"}
	m.Register(a)
	m.Register(b)
	got := m.Components()
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatal("components not in registration order")
	}
}

func TestInitOrder(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec})
	m.Register(&fakeComponent{name: "b", rec: rec})
	if err := m.Init(context.Background()); err != nil {
		t.Fatalf("init: %v", err)
	}
	if got := rec.snapshot(); !equal(got, []string{"a.init", "b.init"}) {
		t.Fatalf("init order: %v", got)
	}
}

func TestInitFailureRollback(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec})
	m.Register(&fakeComponent{name: "b", rec: rec, initErr: errors.New("initfail")})
	if err := m.Init(context.Background()); err == nil {
		t.Fatal("expected init error")
	}
	if got := rec.snapshot(); !equal(got, []string{"a.init", "b.init", "a.stop"}) {
		t.Fatalf("init rollback: %v", got)
	}
	if m.Phase() != PhaseFailed {
		t.Fatalf("phase=%v", m.Phase())
	}
}

func TestStopPropagatesError(t *testing.T) {
	rec := &recorder{}
	m := NewManager()
	m.Register(&fakeComponent{name: "a", rec: rec, stopErr: errors.New("stopfail")})
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := m.Stop(context.Background()); err == nil {
		t.Fatal("expected stop error")
	}
}

func TestEmptyManagerRun(t *testing.T) {
	m := NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := m.Run(ctx); err != nil {
		t.Fatalf("run: %v", err)
	}
}
