// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrRuntime is the base error class for all errors occurring during
	// runtime (and not during the parsing of a scenario or spec)
	ErrRuntime = errors.New("runtime error")
	// ErrRequiredFixture is returned when a required fixture has not
	// been registered with the context.
	ErrRequiredFixture = fmt.Errorf(
		"%w: required fixture missing",
		ErrRuntime,
	)
	// ErrTimeout is returned when a context deadline was exceeeded, a signal
	// was killed in an exec.Spec.Run() call or an expected test spec did not
	// complete in some allocated amount of time.
	ErrTimeout = fmt.Errorf(
		"%w: timeout exceeded",
		ErrRuntime,
	)
)

// RequiredFixtureMissing returns an ErrRequiredFixture with the supplied
// fixture name
func RequiredFixtureMissing(name string) error {
	return fmt.Errorf("%w: %s", ErrRequiredFixture, name)
}

// RuntimeErrors is a collection of zero or more errors resulting from Run()
// calls. It implements the error interface.
type RuntimeErrors struct {
	errors []error
}

func (r *RuntimeErrors) AppendIf(err error) {
	if err != nil {
		r.errors = append(r.errors, err)
	}
}

func (r *RuntimeErrors) Error() string {
	var b strings.Builder
	for x, e := range r.errors {
		b.WriteString(strconv.Itoa(x))
		b.WriteString(": ")
		b.WriteString(e.Error())
		b.WriteRune('\n')
	}
	return b.String()
}

func NewRuntimeErrors() *RuntimeErrors {
	return &RuntimeErrors{
		errors: []error{},
	}
}
