// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package run

import "github.com/gdt-dev/core/testunit"

type Option func(*Run)

// New returns a new Run object that stores test run state.
func New(opts ...Option) *Run {
	r := &Run{
		scenarioResults: map[string][]testunit.Result{},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
