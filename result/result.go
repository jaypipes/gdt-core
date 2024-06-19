// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package result

import (
	"errors"
	"fmt"

	gdterrors "github.com/gdt-dev/gdt/errors"
)

// Result is returned from a `Evaluable.Eval` execution. It serves two
// purposes:
//
// 1) to return an error, if any, from the Eval execution. This error will
// *always* be a `gdterrors.RuntimeError`. Failed assertions are not errors.
// 2) to pass back information about the Run that can be injected into the
// context's `PriorRun` cache. Some plugins, e.g. the gdt-http plugin, use
// cached data from a previous run in order to construct current Run fixtures.
// In the case of the gdt=http plugin, the previous `nethttp.Response` is
// returned in the Result and the `Scenario.Run` method injects that
// information into the context that is supplied to the next Spec's `Run`.
type Result struct {
	// err is any error that was returned from the Evaluable's execution. This
	// is guaranteed to be a `gdterrors.RuntimeError`.
	err error
	// failures is the collection of error messages from assertion failures
	// that occurred during Eval(). These are *not* `gdterrors.RuntimeError`.
	failures []error
	// data is a map, keyed by plugin name, of data about the spec run. Plugins
	// can place anything they want in here and grab it from the context with
	// the `gdtcontext.RunData()` function. Plugins can either set run data
	// into the context using `gdtcontext.SetRunData()` or use the
	// `result.WithData()` modifier when returning a Result.
	data map[string]interface{}
	// vars is a map, keyed by variable name, of variables the test user saved
	// from an individual test step's output. Plugins can place anything they
	// want in here and grab it from the context with the the
	// `gdtcontext.RunVars()` function. Plugins can either set run data into
	// the context using `gdtcontext.SetRunData()` or use the
	// `result.WithVar()` modifier when returning a Result.  clearing and
	// setting run variables.
	vars map[string]interface{}
}

// HasRuntimeError returns true if the Eval() returned a runtime error, false
// otherwise.
func (r *Result) HasRuntimeError() bool {
	return r.err != nil
}

// RuntimeError returns the runtime error
func (r *Result) RuntimeError() error {
	return r.err
}

// HasData returns true if any of the run data has been set, false otherwise.
func (r *Result) HasData() bool {
	return r.data != nil
}

// HasVars returns true if any of the run variables has been set, false
// otherwise.
func (r *Result) HasVars() bool {
	return r.vars != nil
}

// Data returns the raw run data saved in the result
func (r *Result) Data() map[string]interface{} {
	return r.data
}

// Vars returns the raw run vars saved in the result
func (r *Result) Vars() map[string]interface{} {
	return r.vars
}

// Failed returns true if any assertion failed during Eval(), false otherwise.
func (r *Result) Failed() bool {
	return len(r.failures) > 0
}

// Failures returns the collection of assertion failures that occurred during
// Eval().
func (r *Result) Failures() []error {
	return r.failures
}

// SetData sets a value in the result's run data cache.
func (r *Result) SetData(
	key string,
	val interface{},
) {
	if r.data == nil {
		r.data = map[string]interface{}{}
	}
	r.data[key] = val
}

// SetVars sets a value in the result's run variables cache.
func (r *Result) SetVars(
	key string,
	val interface{},
) {
	if r.vars == nil {
		r.vars = map[string]interface{}{}
	}
	r.vars[key] = val
}

// SetFailures sets the result's collection of assertion failures.
func (r *Result) SetFailures(failures ...error) {
	r.failures = failures
}

type ResultModifier func(*Result)

// WithRuntimeError modifies the Result with the supplied error
func WithRuntimeError(err error) ResultModifier {
	if !errors.Is(err, gdterrors.RuntimeError) {
		msg := fmt.Sprintf("expected %s to be a gdterrors.RuntimeError", err)
		// panic here because a plugin author incorrectly implemented their
		// plugin Spec's Eval() method...
		panic(msg)
	}
	return func(r *Result) {
		r.err = err
	}
}

// WithData modifies the Result with the supplied run data key and value
func WithData(key string, val interface{}) ResultModifier {
	return func(r *Result) {
		r.SetData(key, val)
	}
}

// WithVar modifies the Result with the supplied run variable key and value
func WithVar(key string, val interface{}) ResultModifier {
	return func(r *Result) {
		r.SetVars(key, val)
	}
}

// WithFailures modifies the Result the supplied collection of assertion
// failures
func WithFailures(failures ...error) ResultModifier {
	return func(r *Result) {
		r.SetFailures(failures...)
	}
}

// New returns a new Result
func New(mods ...ResultModifier) *Result {
	r := &Result{}
	for _, mod := range mods {
		mod(r)
	}
	return r
}
