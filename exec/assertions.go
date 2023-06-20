// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PipeAssertions contains assertions about the contents of a pipe
type PipeAssertions struct {
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

// Assert checks all the assertions in the PipeAssertions against the supplied
// pipe contents
func (a *PipeAssertions) Assert(
	t *testing.T,
	pipeName string,
	pipe io.Reader,
) {
	assert := assert.New(t)
	buf := &bytes.Buffer{}
	buf.ReadFrom(pipe)

	contents := strings.TrimSpace(buf.String())
	if a.Is != nil {
		assert.Equal(*a.Is, contents)
	}
	if len(a.Contains) > 0 {
		for _, find := range a.Contains {
			assert.Contains(contents, find)
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
			assert.Fail(
				"expected to find one of %s in %s.",
				a.ContainsOneOf, pipeName,
			)
		}
	}
}
