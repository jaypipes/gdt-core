// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	// BaseFields is a list of the field names that specializations of the Spec
	// struct don't need to parse. These are the lower_snake_cased versions
	// that are expected in the YAML definitions.
	BaseFields = []string{
		"name",
		"description",
		"timeout",
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
	// Timeout is the amount of time that the test unit should complete within.
	// Specify a duration using Go's time duration string.
	// See https://pkg.go.dev/time#ParseDuration
	Timeout string `yaml:"timeout,omitempty"`
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

// HasTimeout returns true if the spec has a non-zero timeout duration
func (s *Spec) HasTimeout() bool {
	return s.Timeout != ""
}

// TimeoutDuration returns the duration the spec should execute before timing
// out
func (s *Spec) TimeoutDuration() time.Duration {
	// Parsing already validated the timeout string so no need to check again
	// here
	dur, _ := time.ParseDuration(s.Timeout)
	return dur
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
