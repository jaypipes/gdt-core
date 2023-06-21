// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gdtcontext "github.com/jaypipes/gdt-core/context"
	"github.com/jaypipes/gdt-core/errors"
	gdterrors "github.com/jaypipes/gdt-core/errors"
	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/result"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func (s *fooSpec) Run(ctx context.Context, t *testing.T) error {
	t.Run(s.Title(), func(t *testing.T) {
		assert := assert.New(t)
		// This is just a silly test to demonstrate how to write Run() commands
		// for plugin Spec specialization classes.
		if s.Name == "bar" {
			assert.Equal(s.Foo, "bar")
		} else {
			assert.Equal(s.Foo, "baz")
		}
	})
	return nil
}

func TestRun(t *testing.T) {
	require := require.New(t)
	reg := plugin.NewRegistry()

	reg.Add(&fooPlugin{})

	fp := filepath.Join("testdata", "foo.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	ctx := gdtcontext.New(
		gdtcontext.WithPlugins(
			reg.List(),
		),
	)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
		scenario.WithContext(ctx),
	)
	require.Nil(err)
	require.NotNil(s)

	s.Run(ctx, t)
}

const priorRunDataKey = "priorrun"

type priorRunDefaults struct{}

func (d *priorRunDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type priorRunSpec struct {
	spec.Spec
	State string `yaml:"state"`
	Prior string `yaml:"prior"`
}

func (s *priorRunSpec) UnmarshalYAML(node *yaml.Node) error {
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
		case "state":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.State = valNode.Value
		case "prior":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Prior = valNode.Value
		default:
			if lo.Contains(spec.BaseFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func (s *priorRunSpec) Run(ctx context.Context, t *testing.T) error {
	t.Run(s.Title(), func(t *testing.T) {
		assert := assert.New(t)
		// Here we test that the prior run data that we save at the end of each
		// Run() is showing up properly in the next Run()'s context.
		prData := gdtcontext.PriorRun(ctx)
		if s.Index == 0 {
			assert.Empty(prData)
		} else {
			assert.Contains(prData, priorRunDataKey)
			assert.IsType(prData[priorRunDataKey], "")
			assert.Equal(s.Prior, prData[priorRunDataKey])
		}
	})
	return result.New(result.WithData(priorRunDataKey, s.State))
}

type priorRunPlugin struct{}

func (p *priorRunPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "priorRun",
	}
}

func (p *priorRunPlugin) Defaults() yaml.Unmarshaler {
	return &priorRunDefaults{}
}

func (p *priorRunPlugin) Specs() []gdttypes.Spec {
	return []gdttypes.Spec{&priorRunSpec{}}
}

func TestPriorRun(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "prior-run.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	ctx := gdtcontext.New()
	ctx = gdtcontext.RegisterPlugin(ctx, &priorRunPlugin{})

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
		scenario.WithContext(ctx),
	)
	require.Nil(err)
	require.NotNil(s)

	s.Run(ctx, t)
}

func TestMissingFixtures(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	reg := plugin.NewRegistry()

	reg.Add(&fooPlugin{})

	fp := filepath.Join("testdata", "foo-fixtures.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	ctx := gdtcontext.New(
		gdtcontext.WithPlugins(
			reg.List(),
		),
	)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
		scenario.WithContext(ctx),
	)
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(ctx, t)
	assert.NotNil(err)
	assert.ErrorIs(err, gdterrors.ErrRuntime)
	assert.ErrorIs(err, gdterrors.ErrRequiredFixture)
}
