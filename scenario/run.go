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
	"github.com/jaypipes/gdt-core/result"
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
			err := spec.Run(ctx, t)
			if res, ok := err.(*result.Result); ok {
				// Results can have arbitrary run data stored in them and we
				// save this prior run data in the context (and pass that
				// context to the next Run invocation).
				if res.HasData() {
					ctx = gdtcontext.StorePriorRun(ctx, res.Data())
				}
				errs.AppendIf(res.Unwrap())
			} else {
				errs.AppendIf(err)
			}
		}
	})
	return errs
}
