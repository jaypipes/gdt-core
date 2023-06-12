// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	gopath "path"
	"testing"

	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
)

// Scenario is a generalized gdt test case file. It contains a set of Runnable
// test units.
type Scenario struct {
	// Path is the filepath to the test case.
	Path string `json:"-"`
	// Name is the short name for the test case. If empty, defaults to the base
	// filename in Path.
	Name string `json:"name,omitempty"`
	// Description is a description of the tests contained in the test case.
	Description string `json:"description,omitempty"`
	// Defaults contains any default configuration values for test specs
	// contained within the test scenario.
	//
	// During parsing, plugins are handed this raw data and asked to interpret
	// it into known configuration values for that plugin.
	Defaults map[string]interface{} `json:"defaults,omitempty"`
	// Require specifies an ordered list of fixtures the test case depends on.
	Require []string `json:"require,omitempty"`
	// Tests is the collection of test units in this test case.
	Tests []*spec.Spec `json:"tests,omitempty"`
	// units is a collection of tests that are run as part of this file. These
	// will be the fully parsed and materialized plugin TestSpec structs.
	units []gdttypes.Runnable `json:"-"`
}

// ScenarioModifier sets some value on the test suite
type ScenarioModifier func(s *Scenario)

// WithName sets a test suite's Name attribute
func WithName(name string) ScenarioModifier {
	return func(s *Scenario) {
		s.Name = name
	}
}

// WithPath sets a test suite's Path attribute
func WithPath(path string) ScenarioModifier {
	return func(s *Scenario) {
		s.Path = path
	}
}

// WithDescription sets a test suite's Description attribute
func WithDescription(description string) ScenarioModifier {
	return func(s *Scenario) {
		s.Description = description
	}
}

// WithDefaults sets a test suite's Defaults attribute
func WithDefaults(defaults map[string]interface{}) ScenarioModifier {
	return func(s *Scenario) {
		s.Defaults = defaults
	}
}

// WithRequires sets a test suite's Requires attribute
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

// Append appends a runnable test element to the test case
func (s *Scenario) Append(r gdttypes.Runnable) {
	s.units = append(s.units, r)
}

// Run executes the tests in the test case
func (s *Scenario) Run(ctx context.Context, t *testing.T) {
	t.Run(s.Title(), func(t *testing.T) {
		for _, unit := range s.units {
			unit.Run(ctx, t)
		}
	})
}
