// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin

import (
	"strings"
	"sync"

	gdttypes "github.com/jaypipes/gdt-core/types"
)

// PluginRegistry stores a set of Plugins and is safe to use in threaded
// environments.
type PluginRegistry struct {
	sync.RWMutex
	entries map[string]gdttypes.Plugin
}

// Remove delists the Plugin with registry. Only really useful for testing.
func (r *PluginRegistry) Remove(p gdttypes.Plugin) {
	r.Lock()
	defer r.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	delete(r.entries, lowered)
}

// Add registers a Plugin with the registry.
func (r *PluginRegistry) Add(p gdttypes.Plugin) {
	r.Lock()
	defer r.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	r.entries[lowered] = p
}

// List returns a slice of Plugins that are registered with gdt.
func (r *PluginRegistry) List() []gdttypes.Plugin {
	r.RLock()
	defer r.RUnlock()
	res := []gdttypes.Plugin{}
	for _, p := range r.entries {
		res = append(res, p)
	}
	return res
}

// NewRegistry returns a new PluginRegistry
func NewRegistry() *PluginRegistry {
	return &PluginRegistry{
		entries: map[string]gdttypes.Plugin{},
	}
}
