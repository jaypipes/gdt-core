// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jaypipes/gdt-core/scenario"
	"github.com/stretchr/testify/require"
)

func TestNoExitCodeSimpleCommand(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestExitCode(t *testing.T) {
	require := require.New(t)

	fname := "ls-with-exit-code.yaml"
	// Yay, different exit codes for the same not found error...
	if runtime.GOOS == "darwin" {
		fname = "mac-ls-with-exit-code.yaml"
	}

	fp := filepath.Join("testdata", fname)
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestShellList(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "shell-ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestIs(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "echo-cat.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestContains(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-contains.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestContainsOneOf(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-contains-one-of.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}

func TestSleepTimeout(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "sleep-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}
