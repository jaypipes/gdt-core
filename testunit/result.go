// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package testunit

import (
	"time"

	"github.com/gdt-dev/core/api"
)

func NewResult(
	index int,
	tu *TestUnit,
	res *api.Result,
) Result {
	return Result{
		index:    index,
		name:     tu.Name(),
		elapsed:  tu.Elapsed(),
		skipped:  tu.Skipped(),
		failures: res.Failures(),
		detail:   tu.Detail(),
	}
}

// Result stores a summary of the test execution of a single test unit.
type Result struct {
	// index is the 0-based index of the test unit within the test scenario.
	index int
	// name is the short name of the test unit
	name string
	// skipped is true if the test unit was skipped
	skipped bool
	// failures is the collection of assertion failures for the test spec that
	// occurred during the run. this will NOT include RuntimeErrors.
	failures []error
	// elapsed is the time take to execute the test unit
	elapsed time.Duration
	// detail is a buffer holding any log entries made during the run of the
	// test spec.
	detail string
}

// OK returns whether the test unit executed without any failed assertions.
func (r Result) OK() bool {
	return len(r.failures) == 0
}

// Name returns the name of the test unit.
func (r Result) Name() string {
	return r.name
}

// Index returns the 0-based index of the test unit within its containing test
// scenario.
func (r Result) Index() int {
	return r.index
}

// Failures returns a slice of error messages indicating the test unit
// assertion failures.
func (r Result) Failures() []error {
	return r.failures
}

// Failed returns whether the test unit failed any test assertions.
func (r Result) Failed() bool {
	return !r.skipped && len(r.failures) > 0
}

// Skipped returns whether the test unit was skipped during execution.
func (r Result) Skipped() bool {
	return r.skipped
}

// Detail returns the collected details/output of the test unit.
func (r Result) Detail() string {
	return r.detail
}

// Elapsed returns the elapsed time of the test unit.
func (r Result) Elapsed() time.Duration {
	return r.elapsed
}

// XUnit returns the testunit.Result as a struct that can be serialized to
// JUnit/XUnit XML or JSON.
func (r Result) XUnit() api.XUnitTestCase {
	tc := api.XUnitTestCase{
		Name: r.Name(),
	}
	if r.Skipped() {
		tc.Skipped = true
		tc.Status = "skip"
	} else if r.OK() {
		tc.Status = "ok"
	} else {
		tc.Status = "fail"
	}
	tc.Time = r.Elapsed().String()
	tc.SystemOut = &api.XUnitTextBlock{
		Text: r.Detail(),
	}
	return tc
}
