// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrFailure is the base error class for all errors that represent failed
	// assertions when evaluating a test.
	ErrFailure = errors.New("assertion failed")
)

// TimeoutExceeded returns an ErrFailure when a test's execution exceeds a
// timeout length.
func TimeoutExceeded(duration string) error {
	return fmt.Errorf(
		"%s: timeout of %s exceeded",
		ErrFailure, duration,
	)
}

// NotEqualLength returns an ErrFailure when an expected length doesn't equal an
// observed length.
func NotEqualLength(exp, got int) error {
	return fmt.Errorf(
		"%w: expected length of %d but got %d",
		ErrFailure, exp, got,
	)
}

// NotEqual returns an ErrFailure when an expected thing doesn't equal an
// observed thing.
func NotEqual(exp, got interface{}) error {
	return fmt.Errorf(
		"%w: expected %v but got %v",
		ErrFailure, exp, got,
	)
}

// NotIn returns an ErrFailure when an expected thing doesn't appear in an
// expected container.
func NotIn(element, container interface{}) error {
	return fmt.Errorf(
		"%w: expected %v to contain %v",
		ErrFailure, container, element,
	)
}

// NoneIn returns an ErrFailure when none of a list of elements appears in an
// expected container.
func NoneIn(elements, container interface{}) error {
	return fmt.Errorf(
		"%w: expected %v to contain one of %v",
		ErrFailure, container, elements,
	)
}
