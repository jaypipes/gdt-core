// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/google/shlex"

	gdtcontext "github.com/jaypipes/gdt-core/context"
	gdterrors "github.com/jaypipes/gdt-core/errors"
)

// Run executes the specific exec test spec.
func (s *Spec) Run(ctx context.Context, t *testing.T) error {

	assertions := newAssertions(
		s.ExitCode, s.Out, s.Err,
	)
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}

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
	outbuf.ReadFrom(outpipe)
	errbuf.ReadFrom(errpipe)

	err = cmd.Wait()
	if gdtcontext.TimedOut(ctx, err) {
		return gdterrors.ErrTimeout
	}
	ec := 0
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		ec = eerr.ExitCode()
	}

	if !assertions.OK(ec, outbuf, errbuf) {
		for _, failure := range assertions.Failures() {
			t.Error(failure)
		}
	}
	return nil
}
