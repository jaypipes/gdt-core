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
	Name string `json:"name,omitempty"`
	// Description of the test unit
	Description string `json:"description,omitempty"`
}

// SpecModifier sets some value on the test suite
type SpecModifier func(s *Spec)

// WithName sets a test suite's Name attribute
func WithName(name string) SpecModifier {
	return func(s *Spec) {
		s.Name = name
	}
}

// WithDescription sets a test suite's Description attribute
func WithDescription(description string) SpecModifier {
	return func(s *Spec) {
		s.Description = description
	}
}

// New returns a new Spec
func New(mods ...SpecModifier) *Spec {
	s := &Spec{}
	for _, mod := range mods {
		mod(s)
	}
	return s
}

// Run executes the specific test unit
func (s *Spec) Run(ctx context.Context, t *testing.T) {

}
