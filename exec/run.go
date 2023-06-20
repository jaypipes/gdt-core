// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"context"
	"os/exec"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run executes the specific exec test spec.
func (s *ExecSpec) Run(ctx context.Context, t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var cmd *exec.Cmd
	if s.Shell == "" {
		args, err := shlex.Split(s.Exec)
		if err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(s.Shell, "-c", s.Exec)
	}

	outpipe, err := cmd.StdoutPipe()
	errpipe, err := cmd.StderrPipe()
	require.Nil(err)

	err = cmd.Start()
	require.Nil(err)

	if s.Out != nil {
		s.Out.Assert(t, "stdout", outpipe)
	}
	if s.Err != nil {
		s.Err.Assert(t, "stderr", errpipe)
	}

	err = cmd.Wait()
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		assert.Equal(s.ExitCode, eerr.ExitCode())
	}
}
