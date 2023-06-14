// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	"fmt"
	gopath "path"
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidScenarioDefinition = fmt.Errorf("Invalid scenario YAML definition.")
)

// Scenario is a generalized gdt test case file. It contains a set of Runnable
// test units.
type Scenario struct {
	// plugins is the list of plugins known to the scenario. This is injected
	// during scenario file parsing.
	plugins []plugin.Plugin
	// Path is the filepath to the test case.
	Path string `yaml:"-"`
	// Name is the short name for the test case. If empty, defaults to the base
	// filename in Path.
	Name string `yaml:"name,omitempty"`
	// Description is a description of the tests contained in the test case.
	Description string `yaml:"description,omitempty"`
	// Defaults contains any default configuration values for test specs
	// contained within the test scenario.
	//
	// During parsing, plugins are handed this raw data and asked to interpret
	// it into known configuration values for that plugin.
	Defaults map[string]interface{} `yaml:"defaults,omitempty"`
	// Require specifies an ordered list of fixtures the test case depends on.
	Require []string `yaml:"require,omitempty"`
	// Tests is the collection of test units in this test case. These will be
	// the fully parsed and materialized plugin Spec structs.
	Tests []gdttypes.Runnable `yaml:"tests,omitempty"`
}

// ScenarioModifier sets some value on the test scenario
type ScenarioModifier func(s *Scenario)

// WithName sets a test scenario's Name attribute
func WithName(name string) ScenarioModifier {
	return func(s *Scenario) {
		s.Name = name
	}
}

// WithPath sets a test scenario's Path attribute
func WithPath(path string) ScenarioModifier {
	return func(s *Scenario) {
		s.Path = path
	}
}

// WithDescription sets a test scenario's Description attribute
func WithDescription(description string) ScenarioModifier {
	return func(s *Scenario) {
		s.Description = description
	}
}

// WithDefaults sets a test scenario's Defaults attribute
func WithDefaults(defaults map[string]interface{}) ScenarioModifier {
	return func(s *Scenario) {
		s.Defaults = defaults
	}
}

// WithRequires sets a test scenario's Requires attribute
func WithRequires(require []string) ScenarioModifier {
	return func(s *Scenario) {
		s.Require = require
	}
}

// New returns a new Scenario
func New(mods ...ScenarioModifier) *Scenario {
	s := &Scenario{}
	for _, mod := range mods {
		mod(s)
	}
	return s
}

// Title returns the Name of the scenario or the Path's file/base name if there
// is no name.
func (s *Scenario) Title() string {
	if s.Name != "" {
		return s.Name
	}
	return gopath.Base(s.Path)
}

// Run executes the tests in the test case
func (s *Scenario) Run(ctx context.Context, t *testing.T) {
	t.Run(s.Title(), func(t *testing.T) {
		for _, spec := range s.Tests {
			spec.Run(ctx, t)
		}
	})
}

// UnmarshalYAML is a custom unmarshaler that asks plugins for their known spec
// types and attempts to unmarshal test spec contents into those types.
func (s *Scenario) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return ErrInvalidScenarioDefinition
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return ErrInvalidScenarioDefinition
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return ErrInvalidScenarioDefinition
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return ErrInvalidScenarioDefinition
			}
			s.Description = valNode.Value
		case "require":
			if valNode.Kind != yaml.SequenceNode {
				return ErrInvalidScenarioDefinition
			}
			requires := make([]string, len(valNode.Content))
			for x, n := range valNode.Content {
				requires[x] = n.Value
			}
			s.Require = requires
		case "defaults":
			if valNode.Kind != yaml.MappingNode {
				return ErrInvalidScenarioDefinition
			}
			defaults := map[string]interface{}{}
			if err := valNode.Decode(&defaults); err != nil {
				return err
			}
			s.Defaults = defaults
		case "tests":
			if valNode.Kind != yaml.SequenceNode {
				return ErrInvalidScenarioDefinition
			}
			for _, testNode := range valNode.Content {
				for _, p := range s.plugins {
					specs := p.Specs()
					for _, spec := range specs {
						if err := testNode.Decode(spec); err != nil {
							return err
						} else {
							s.Tests = append(s.Tests, spec.(gdttypes.Runnable))
						}
					}
				}
			}
		}

	}
	return nil
}
