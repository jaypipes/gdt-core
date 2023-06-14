// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec

import (
	"context"
	"testing"
)

// Spec represents a single test action and one or more assertions about
// output or behaviour. All gdt plugins have their own Spec structs that
// inherit from this base struct.
type Spec struct {
	// Name for the individual test unit
	Name string `yaml:"name,omitempty"`
	// Description of the test unit
	Description string `yaml:"description,omitempty"`
}

// Run executes the specific test unit.
//
// NOTE(jaypipes): consider this a pure virtual function. Should not be
// executed since the plugin-specific subclass should implement its own Run()
// method.
func (s *Spec) Run(ctx context.Context, t *testing.T) {}
