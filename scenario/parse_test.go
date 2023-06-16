// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jaypipes/gdt-core/errors"
	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type failDefaults struct{}

func (d *failDefaults) UnmarshalYAML(node *yaml.Node) error {
	return fmt.Errorf("Indy, bad dates!")
}

type failSpec struct {
	spec.Spec
}

func (s *failSpec) Run(context.Context, *testing.T) {}

func (s *failSpec) UnmarshalYAML(node *yaml.Node) error {
	return fmt.Errorf("Indy, bad dates!")
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

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "foo":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Foo = valNode.Value
		case "name", "description":
			continue
		default:
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

type barDefaults struct {
	Foo string `yaml:"bar"`
}

func (d *barDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type barSpec struct {
	spec.Spec
	Bar int `yaml:"bar"`
}

func (s *barSpec) Run(context.Context, *testing.T) {}

func (s *barSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "bar":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			if v, err := strconv.Atoi(valNode.Value); err != nil {
				return errors.ExpectedIntAt(valNode)
			} else {
				s.Bar = v
			}
		case "name", "description":
			continue
		default:
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func TestNoPlugins(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.ErrorIs(err, errors.ErrUnknownSpec)
	assert.Nil(s)
}

func TestNoTests(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "no-tests.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	// When there are no plugins and no tests, we should successfully parse the
	// scenario's defaults and have an empty set of Tests in the scenario
	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	sc := s.(*scenario.Scenario)
	assert.Equal("no-tests", sc.Name)
	assert.Equal(filepath.Join("testdata", "no-tests.yaml"), sc.Path)
	assert.Equal([]string{"books_api", "books_data"}, sc.Require)
	assert.Equal(
		map[string]interface{}{
			"http": map[string]interface{}{
				"base_url": "http://127.0.0.1:4000",
			},
		},
		sc.Defaults,
	)
	assert.Empty(sc.Tests)
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

func (p *failingPlugin) Specs() []gdttypes.Spec {
	return []gdttypes.Spec{&failSpec{}}
}

func TestFailingPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	p := &failingPlugin{}

	plugin.Register(p)
	defer plugin.Unregister(p)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.NotErrorIs(err, errors.ErrUnknownSpec)
	assert.Nil(s)
}

type barPlugin struct{}

func (p *barPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "bar",
	}
}

func (p *barPlugin) Defaults() yaml.Unmarshaler {
	return &barDefaults{}
}

func (p *barPlugin) Specs() []gdttypes.Spec {
	return []gdttypes.Spec{&barSpec{}}
}

func TestUnknownSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	p := &barPlugin{}

	plugin.Register(p)
	defer plugin.Unregister(p)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.Nil(s)
}

type fooPlugin struct{}

func (p *fooPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "foo",
	}
}

func (p *fooPlugin) Defaults() yaml.Unmarshaler {
	return &fooDefaults{}
}

func (p *fooPlugin) Specs() []gdttypes.Spec {
	return []gdttypes.Spec{&fooSpec{}}
}

func TestKnownSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	p := &fooPlugin{}

	plugin.Register(p)
	defer plugin.Unregister(p)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	sc := s.(*scenario.Scenario)
	assert.Equal("foo-test", sc.Name)
	assert.Equal(filepath.Join("testdata", "foo-test.yaml"), sc.Path)
	assert.Empty(sc.Require)
	assert.Equal(
		map[string]interface{}{
			"foo": map[string]interface{}{
				"key": "value",
			},
		},
		sc.Defaults,
	)
	expTests := []gdttypes.Runnable{
		&fooSpec{
			Spec: spec.Spec{
				Index: 0,
				Name:  "bar",
			},
			Foo: "bar",
		},
		&fooSpec{
			Spec: spec.Spec{
				Index:       1,
				Description: "Bazzy Bizzy",
			},
			Foo: "baz",
		},
	}
	assert.Equal(expTests, sc.Tests)
}
