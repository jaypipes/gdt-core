// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *fooSpec) Run(ctx context.Context, t *testing.T) {
	t.Run(s.Title(), func(t *testing.T) {
		assert := assert.New(t)
		// This is just a silly test to demonstrate how to write Run() commands
		// for plugin Spec specialization classes.
		if s.Name == "bar" {
			assert.Equal(s.Foo, "bar")
		} else {
			assert.Equal(s.Foo, "baz")
		}
	})
}

func TestRun(t *testing.T) {
	require := require.New(t)
	reg := plugin.NewRegistry()

	reg.Add(&fooPlugin{})

	fp := filepath.Join("testdata", "foo.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
		scenario.WithPlugins(reg.List()),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	s.Run(ctx, t)
}
