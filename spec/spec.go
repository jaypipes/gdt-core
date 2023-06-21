// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec

import (
	"context"
	"strconv"
	"strings"
	"testing"
)

var (
	// BaseFields is a list of the field names that specializations of the Spec
	// struct don't need to parse. These are the lower_snake_cased versions
	// that are expected in the YAML definitions.
	BaseFields = []string{
		"name",
		"description",
	}
)

// Spec represents a single test action and one or more assertions about
// output or behaviour. All gdt plugins have their own Spec structs that
// inherit from this base struct.
type Spec struct {
	// Index within the scenario where this Spec is located
	Index int `yaml:"-"`
	// Name for the individual test unit
	Name string `yaml:"name,omitempty"`
	// Description of the test unit
	Description string `yaml:"description,omitempty"`
}

// Title returns the Name of the scenario or the Path's file/base name if there
// is no name.
func (s *Spec) Title() string {
	if s.Name != "" {
		return s.Name
	}
	if s.Description != "" {
		return slugify(s.Description)
	}
	return strconv.Itoa(s.Index)
}

// slugify returns a new string that lowercases and removes spaces and forward
// slashes from the supplied string
func slugify(s string) string {
	s = strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.TrimSpace(s),
				" ", "-"),
			"/", "-",
		),
	)
	for {
		if strings.Contains(s, "--") {
			s = strings.ReplaceAll(s, "--", "-")
		} else {
			return s
		}
	}
}

// Run executes the specific test unit.
//
// NOTE(jaypipes): consider this a pure virtual function. Should not be
// executed since the plugin-specific subclass should implement its own Run()
// method.
func (s *Spec) Run(ctx context.Context, t *testing.T) error {
	return nil
}
