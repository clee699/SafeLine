package upstream

import (
	"errors"
	"fmt"
)

// Upstream represents the upstream service.
type Upstream struct {
	Name string
	URL  string
}

// Manager manages upstream services.
type Manager struct {
	upstreams map[string]Upstream
}

// NewManager creates a new upstream manager.
func NewManager() *Manager {
	return &Manager{
		upstreams: make(map[string]Upstream),
	}
}

// Add adds a new upstream service.
func (m *Manager) Add(name, url string) error {
	if _, exists := m.upstreams[name]; exists {
		return errors.New("upstream already exists")
	}
	m.upstreams[name] = Upstream{Name: name, URL: url}
	return nil
}

// Remove removes an upstream service.
func (m *Manager) Remove(name string) error {
	if _, exists := m.upstreams[name]; !exists {
		return errors.New("upstream not found")
	}
	delete(m.upstreams, name)
	return nil
}

// List lists all upstream services.
func (m *Manager) List() []Upstream {
	upstreamsList := []Upstream{}
	for _, upstream := range m.upstreams {
		upstreamsList = append(upstreamsList, upstream)
	}
	return upstreamsList
}

// Update updates an existing upstream service.
func (m *Manager) Update(name, url string) error {
	if _, exists := m.upstreams[name]; !exists {
		return errors.New("upstream not found")
	}
	m.upstreams[name] = Upstream{Name: name, URL: url}
	return nil
}

// Show displays the details of an upstream service.
func (m *Manager) Show(name string) (Upstream, error) {
	upstream, exists := m.upstreams[name]
	if !exists {
		return Upstream{}, errors.New("upstream not found")
	}
	return upstream, nil
}