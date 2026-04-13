package adapter

import (
	"fmt"
	"sort"
	"sync"
)

// Registry stores adapters by their stable identity.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]Adapter),
	}
}

func (r *Registry) Register(a Adapter) error {
	if a == nil {
		return fmt.Errorf("adapter.Registry.Register: adapter cannot be nil")
	}

	id := a.Identity()
	if id == "" {
		return fmt.Errorf("adapter.Registry.Register: adapter identity cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adapters[id]; exists {
		return fmt.Errorf("adapter.Registry.Register: adapter already registered: %s", id)
	}

	r.adapters[id] = a
	return nil
}

func (r *Registry) MustRegister(a Adapter) {
	if err := r.Register(a); err != nil {
		panic(err)
	}
}

func (r *Registry) Resolve(id string) (Adapter, error) {
	if id == "" {
		return nil, fmt.Errorf("adapter.Registry.Resolve: adapter identity cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	a, exists := r.adapters[id]
	if !exists {
		return nil, fmt.Errorf("adapter.Registry.Resolve: adapter not registered: %s", id)
	}

	return a, nil
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.adapters))
	for id := range r.adapters {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
