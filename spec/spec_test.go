// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec_test

import (
	"testing"

	"github.com/jaypipes/gdt-core/spec"
	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	assert := assert.New(t)

	s := spec.New(
		spec.WithName("foo"),
	)

	assert.Equal("foo", s.Name)
	assert.Equal("", s.Description)
}
