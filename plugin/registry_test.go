// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin_test

import (
	"context"
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type fooDefaults struct {
	Foo string `yaml:"foo"`
}

func (d *fooDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type fooSpec struct {
	spec.Spec
	Foo string `yaml:"foo"`
}

func (s *fooSpec) Run(context.Context, *testing.T) error {
	return nil
}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type fooPlugin struct{}

func (p *fooPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "foo",
	}
}

func (p *fooPlugin) Defaults() yaml.Unmarshaler {
	return &fooDefaults{}
}

func (p *fooPlugin) Specs() []gdttypes.Spec {
	return []gdttypes.Spec{&fooSpec{}}
}

func TestRegisterAndList(t *testing.T) {
	assert := assert.New(t)

	r := plugin.NewRegistry()

	plugins := r.List()
	assert.Equal(0, len(plugins))

	r.Add(&fooPlugin{})

	plugins = r.List()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)

	// Add called twice with the same named plugin should be be a no-op

	r.Add(&fooPlugin{})

	plugins = r.List()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)
}
