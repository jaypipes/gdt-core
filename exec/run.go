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

	gdtcontext "github.com/jaypipes/gdt-core/context"
	gdterrors "github.com/jaypipes/gdt-core/errors"
)

// Run executes the specific exec test spec.
func (s *ExecSpec) Run(ctx context.Context, t *testing.T) error {
	assert := assert.New(t)

	var cmd *exec.Cmd
	if s.Shell == "" {
		args, err := shlex.Split(s.Exec)
		if err != nil {
			return err
		}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, s.Shell, "-c", s.Exec)
	}

	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if gdtcontext.TimedOut(ctx, err) {
		return gdterrors.ErrTimeout
	}
	if err != nil {
		return err
	}

	if s.Out != nil {
		s.Out.Assert(t, "stdout", outpipe)
	}
	if s.Err != nil {
		s.Err.Assert(t, "stderr", errpipe)
	}

	err = cmd.Wait()
	if gdtcontext.TimedOut(ctx, err) {
		return gdterrors.ErrTimeout
	}
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		assert.Equal(s.ExitCode, eerr.ExitCode())
		return err
	}
	return nil
}
