// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin

import (
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	// plugins is the collection of known gdt plugins. Plugins call the
	// `RegisterPlugin` function to register themselves during `init`.
	plugins = pluginRegistry{
		entries: map[string]Plugin{},
	}
)

// PluginInfo contains basic information about the plugin and what type of
// tests it can handle.
type PluginInfo struct {
	// Name is the primary name of the plugin
	Name string
	// Aliases is an optional set of aliased names for the plugin
	Aliases []string
	// Description describes what types of tests the plugin can handle.
	Description string
}

// Plugin is the driver interface for different types of gdt tests.
type Plugin interface {
	// Info returns a struct that describes what the plugin does
	Info() PluginInfo
	// Defaults returns a YAML Unmarshaler types that the plugin knows how
	// to parse its defaults configuration with.
	Defaults() yaml.Unmarshaler
	// Specs returns a list of YAML Unmarshaler types that the plugin knows
	// how to parse.
	Specs() []yaml.Unmarshaler
}

// pluginRegistry stores all known Plugins
type pluginRegistry struct {
	sync.RWMutex
	entries map[string]Plugin
}

// Unregister delists the Plugin with gdt. Only really useful for testing.
func Unregister(p Plugin) {
	plugins.Lock()
	defer plugins.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	delete(plugins.entries, lowered)
}

// Register registers a Plugin with gdt.
func Register(p Plugin) {
	plugins.Lock()
	defer plugins.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	plugins.entries[lowered] = p
}

// List returns a slice of Plugins that are registered with gdt.
func List() []Plugin {
	plugins.RLock()
	defer plugins.RUnlock()
	res := []Plugin{}
	for _, p := range plugins.entries {
		res = append(res, p)
	}
	return res
}
