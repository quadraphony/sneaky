package adapter

import (
	"context"
	"testing"

	"sneaky-core/internal/runtime"
)

type testAdapter struct {
	id string
}

func (a testAdapter) Identity() string {
	return a.id
}

func (a testAdapter) Capabilities() Capabilities {
	return Capabilities{}
}

func (a testAdapter) ValidateConfig(StartRequest) error {
	return nil
}

func (a testAdapter) Start(context.Context, StartRequest) (runtime.Handle, error) {
	return nil, nil
}

func TestRegistryRegisterResolveAndList(t *testing.T) {
	reg := NewRegistry()

	if err := reg.Register(testAdapter{id: "zeta"}); err != nil {
		t.Fatalf("register zeta: %v", err)
	}
	if err := reg.Register(testAdapter{id: "alpha"}); err != nil {
		t.Fatalf("register alpha: %v", err)
	}

	got, err := reg.Resolve("alpha")
	if err != nil {
		t.Fatalf("resolve alpha: %v", err)
	}
	if got.Identity() != "alpha" {
		t.Fatalf("resolved wrong adapter: got %q", got.Identity())
	}

	ids := reg.List()
	if len(ids) != 2 || ids[0] != "alpha" || ids[1] != "zeta" {
		t.Fatalf("unexpected adapter list: %#v", ids)
	}
}

func TestRegistryRejectsDuplicateIdentity(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(testAdapter{id: "singbox"}); err != nil {
		t.Fatalf("initial register: %v", err)
	}

	if err := reg.Register(testAdapter{id: "singbox"}); err == nil {
		t.Fatal("expected duplicate registration error")
	}
}
