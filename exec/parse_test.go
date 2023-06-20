// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/gdt-core/errors"
	gdtexec "github.com/jaypipes/gdt-core/exec"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnknownShell(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "unknown-shell-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	assert.NotNil(err)
	assert.ErrorIs(err, errors.ErrInvalid)
	assert.Nil(s)
}

func TestSimpleCommand(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	sc := s.(*scenario.Scenario)
	expTests := []gdttypes.Runnable{
		&gdtexec.ExecSpec{
			Spec: spec.Spec{
				Index: 0,
			},
			Exec: "ls",
		},
	}
	assert.Equal(expTests, sc.Tests)
}
