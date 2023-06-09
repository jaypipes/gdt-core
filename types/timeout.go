// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import (
	"time"
)

// Timeout contains information about the duration within which a Spec should
// run along with whether a deadline exceeded/timeout error should be expected
// or not.
type Timeout struct {
	// After is the amount of time that the test unit should complete within.
	// Specify a duration using Go's time duration string.
	// See https://pkg.go.dev/time#ParseDuration
	After string `yaml:"after,omitempty"`
	// Expected indicates whether the timeout is expected to be exceeded. This
	// is mostly useful for unit testing of the timeout functionality itself.
	Expected bool `yaml:"expected,omitempty"`
}

// Duration returns the time duration of the Timeout
func (t *Timeout) Duration() time.Duration {
	// Parsing already validated the timeout string so no need to check again
	// here
	dur, _ := time.ParseDuration(t.After)
	return dur
}
