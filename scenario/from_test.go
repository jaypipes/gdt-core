// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestFromNoPlugins(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "http-failure.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	sc := s.(*scenario.Scenario)
	assert.Equal("HTTP failure", sc.Name)
	assert.Equal("testdata/http-failure.yaml", sc.Path)
	assert.Equal([]string{"books_api", "books_data"}, sc.Require)
	assert.Equal(
		map[string]interface{}{
			"http": map[string]interface{}{
				"base_url": "http://127.0.0.1:4000",
			},
		},
		sc.Defaults,
	)
	// With no plugins, there should be no tests parsed...
	assert.Empty(sc.Tests)
}

type failDefaults struct{}

func (d *failDefaults) UnmarshalYAML(node *yaml.Node) error {
	return fmt.Errorf("Indy, bad dates!")
}

type failSpec struct{}

func (s *failSpec) Run(context.Context, *testing.T) {}

func (s *failSpec) UnmarshalYAML(node *yaml.Node) error {
	return fmt.Errorf("Indy, bad dates!")
}

type failingPlugin struct{}

func (p *failingPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "failer",
	}
}

func (p *failingPlugin) Defaults() yaml.Unmarshaler {
	return &failDefaults{}
}

func (p *failingPlugin) Specs() []yaml.Unmarshaler {
	return []yaml.Unmarshaler{&failSpec{}}
}

func TestFromFailingPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "http-failure.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	p := &failingPlugin{}

	plugin.Register(p)
	defer plugin.Unregister(p)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.Nil(s)
}

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

func (s *fooSpec) Run(context.Context, *testing.T) {}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node")
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return fmt.Errorf("expected scalar key node")
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "foo":
			if valNode.Kind != yaml.ScalarNode {
				return fmt.Errorf("expected foo scalar node")
			}
			s.Name = valNode.Value
		default:
			return fmt.Errorf("unknown field %s", key)
		}
	}
	return nil
}

type unknownSpecPlugin struct{}

func (p *unknownSpecPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "unknownspec",
	}
}

func (p *unknownSpecPlugin) Defaults() yaml.Unmarshaler {
	return &fooDefaults{}
}

func (p *unknownSpecPlugin) Specs() []yaml.Unmarshaler {
	return []yaml.Unmarshaler{&fooSpec{}}
}

func TestFromUnknownSpecPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "http-failure.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	p := &unknownSpecPlugin{}

	plugin.Register(p)
	defer plugin.Unregister(p)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.Nil(s)
}
