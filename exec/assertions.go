// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"strings"

	"github.com/jaypipes/gdt-core/errors"
)

// PipeAssertions contains assertions about the contents of a pipe
type PipeAssertions struct {
	// failures contains the set of error messages for failed assertions
	failures []error
	// terminal indicates there was a failure in evaluating the assertions that
	// should be considered a terminal condition (and therefore the test action
	// should not be retried).
	terminal bool
	// Is contains the exact match (minus whitespace) of the contents of the
	// pipe
	Is *string `yaml:"is,omitempty"`
	// Contains is one or more strings that *all* must be present in the
	// contents of the pipe
	Contains []string `yaml:"contains,omitempty"`
	// ContainsOneOf is one or more strings of which *at least one* must be
	// present in the contents of the pipe
	ContainsOneOf []string `yaml:"contains_one_of,omitempty"`
}

// Fail appends a supplied error to the set of failed assertions
func (a *PipeAssertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of errors for all failed assertions
func (a *PipeAssertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// Terminal returns a bool indicating the assertions failed in a way that is
// not retryable.
func (a *PipeAssertions) Terminal() bool {
	if a == nil {
		return false
	}
	return a.terminal
}

// OK checks all the assertions in the PipeAssertions against the supplied pipe
// contents and returns true if all assertions pass.
func (a *PipeAssertions) OK(
	pipeName string,
	pipe *bytes.Buffer,
) bool {
	if a == nil {
		return true
	}

	res := true
	contents := strings.TrimSpace(pipe.String())
	if a.Is != nil {
		exp := *a.Is
		got := contents
		if exp != got {
			a.Fail(errors.NotEqual(exp, got))
			res = false
		}
	}
	if len(a.Contains) > 0 {
		for _, find := range a.Contains {
			if !strings.Contains(contents, find) {
				a.Fail(errors.NotIn(find, pipeName))
				res = false
			}
		}
	}
	if len(a.ContainsOneOf) > 0 {
		found := false
		for _, find := range a.ContainsOneOf {
			if idx := strings.Index(contents, find); idx > -1 {
				found = true
				break
			}
		}
		if !found {
			a.Fail(errors.NoneIn(a.ContainsOneOf, pipeName))
			res = false
		}
	}
	return res
}

// assertions contains all assertions made for the exec test
type assertions struct {
	// failures contains the set of error messages for failed assertions
	failures []error
	// terminal indicates there was a failure in evaluating the assertions that
	// should be considered a terminal condition (and therefore the test action
	// should not be retried).
	terminal bool
	// exitCode contains the required exit code
	exitCode int
	// outpipe contains the assertions against stdout
	outpipe *PipeAssertions
	// errpipe contains the assertions against stderr
	errpipe *PipeAssertions
}

// Fail appends a supplied error to the set of failed assertions
func (a *assertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of errors for all failed assertions
func (a *assertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// Terminal returns a bool indicating the assertions failed in a way that is
// not retryable.
func (a *assertions) Terminal() bool {
	if a == nil {
		return false
	}
	return a.terminal
}

// OK checks all the assertions against the supplied arguments and returns true
// if all assertions pass.
func (a *assertions) OK(
	exitCode int,
	stdout *bytes.Buffer,
	stderr *bytes.Buffer,
) bool {
	res := true
	if exitCode != a.exitCode {
		a.Fail(errors.NotEqual(exitCode, a.exitCode))
		res = false
	}
	if !a.outpipe.OK("stdout", stdout) {
		a.failures = append(a.failures, a.outpipe.Failures()...)
		res = false
	}
	if !a.errpipe.OK("stderr", stderr) {
		a.failures = append(a.failures, a.errpipe.Failures()...)
		res = false
	}
	return res
}

// newAssertions returns an assertions object populated with the supplied exec
// spec assertions
func newAssertions(
	exitCode int,
	outpipe *PipeAssertions,
	errpipe *PipeAssertions,
) *assertions {
	return &assertions{
		failures: []error{},
		exitCode: exitCode,
		outpipe:  outpipe,
		errpipe:  errpipe,
	}
}
