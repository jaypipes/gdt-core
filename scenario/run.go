// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	"testing"
)

// Run executes the tests in the test scenario
func (s *Scenario) Run(ctx context.Context, t *testing.T) {
	t.Run(s.Title(), func(t *testing.T) {
		for _, spec := range s.Tests {
			spec.Run(ctx, t)
		}
	})
}
