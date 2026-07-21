package di

import (
	"errors"
	"reflect"
	"testing"
)

type testService struct{ val int }

func TestRegisterAndResolveSingleton(t *testing.T) {
	c := New()
	typ := reflect.TypeOf(&testService{})
	if err := c.Register(typ, func(*Container) (any, error) {
		return &testService{val: 1}, nil
	}, Singleton); err != nil {
		t.Fatalf("register: %v", err)
	}
	a, err := c.Resolve(typ)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	b, err := c.Resolve(typ)
	if err != nil {
		t.Fatalf("resolve2: %v", err)
	}
	if a != b {
		t.Fatal("expected singleton to return same instance")
	}
	if a.(*testService).val != 1 {
		t.Fatal("unexpected value")
	}
}

func TestTransient(t *testing.T) {
	c := New()
	typ := reflect.TypeOf(&testService{})
	c.MustRegister(typ, func(*Container) (any, error) {
		return &testService{}, nil
	}, Transient)
	a, _ := c.Resolve(typ)
	b, _ := c.Resolve(typ)
	if a == b {
		t.Fatal("expected transient to return different instances")
	}
}

func TestNotRegistered(t *testing.T) {
	c := New()
	_, err := c.Resolve(reflect.TypeOf(&testService{}))
	if !errors.Is(err, ErrNotRegistered) {
		t.Fatalf("expected ErrNotRegistered, got %v", err)
	}
}

func TestAlreadyRegistered(t *testing.T) {
	c := New()
	typ := reflect.TypeOf(&testService{})
	c.MustRegister(typ, nil, Singleton)
	err := c.Register(typ, nil, Singleton)
	if !errors.Is(err, ErrAlreadyRegistered) {
		t.Fatalf("expected ErrAlreadyRegistered, got %v", err)
	}
}

func TestFactoryError(t *testing.T) {
	c := New()
	typ := reflect.TypeOf(&testService{})
	c.MustRegister(typ, func(*Container) (any, error) {
		return nil, errors.New("boom")
	}, Singleton)
	_, err := c.Resolve(typ)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected boom error, got %v", err)
	}
}

func TestFactoryWrongType(t *testing.T) {
	c := New()
	typ := reflect.TypeOf(&testService{})
	c.MustRegister(typ, func(*Container) (any, error) {
		return 42, nil
	}, Singleton)
	if _, err := c.Resolve(typ); err == nil {
		t.Fatal("expected type mismatch error")
	}
}

func TestCircularDependency(t *testing.T) {
	c := New()
	aType := reflect.TypeOf(&struct{ A int }{})
	bType := reflect.TypeOf(&struct{ B int }{})
	c.MustRegister(aType, func(c *Container) (any, error) {
		if _, err := c.Resolve(bType); err != nil {
			return nil, err
		}
		return &struct{ A int }{}, nil
	}, Singleton)
	c.MustRegister(bType, func(c *Container) (any, error) {
		if _, err := c.Resolve(aType); err != nil {
			return nil, err
		}
		return &struct{ B int }{}, nil
	}, Singleton)
	_, err := c.Resolve(aType)
	if !errors.Is(err, ErrCircularDependency) {
		t.Fatalf("expected ErrCircularDependency, got %v", err)
	}
}

func TestMustRegisterPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	c := New()
	typ := reflect.TypeOf(&testService{})
	c.MustRegister(typ, nil, Singleton)
	c.MustRegister(typ, nil, Singleton)
}

func TestResolveDependency(t *testing.T) {
	c := New()
	bType := reflect.TypeOf(&struct{ B int }{})
	aType := reflect.TypeOf(&struct{ A int }{})
	c.MustRegister(bType, func(*Container) (any, error) {
		return &struct{ B int }{B: 2}, nil
	}, Singleton)
	c.MustRegister(aType, func(c *Container) (any, error) {
		b, err := c.Resolve(bType)
		if err != nil {
			return nil, err
		}
		return &struct{ A int }{A: b.(*struct{ B int }).B * 10}, nil
	}, Singleton)
	v, err := c.Resolve(aType)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if v.(*struct{ A int }).A != 20 {
		t.Fatalf("unexpected A: %d", v.(*struct{ A int }).A)
	}
}
