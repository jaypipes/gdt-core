// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	"strings"
	"testing"

	gdtcontext "github.com/jaypipes/gdt-core/context"
	gdterrors "github.com/jaypipes/gdt-core/errors"
)

// Run executes the tests in the test scenario
func (s *Scenario) Run(ctx context.Context, t *testing.T) error {
	if len(s.Require) > 0 {
		fixtures := gdtcontext.Fixtures(ctx)
		for _, fname := range s.Require {
			lookup := strings.ToLower(fname)
			fix, found := fixtures[lookup]
			if !found {
				return gdterrors.RequiredFixtureMissing(fname)
			}
			fix.Start()
			defer fix.Stop()
		}
	}
	errs := gdterrors.NewRuntimeErrors()
	t.Run(s.Title(), func(t *testing.T) {
		for _, spec := range s.Tests {
			errs.AppendIf(spec.Run(ctx, t))
		}
	})
	return errs
}
