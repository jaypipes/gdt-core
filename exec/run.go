// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

// Run executes the specific exec test spec.
func (s *ExecSpec) Run(ctx context.Context, t *testing.T) {
	assert := assert.New(t)
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
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		assert.Equal(s.ExitCode, eerr.ExitCode())
	}
	if s.Out != nil {
		// evaluate the pipe assertions on stdout
		outContents := strings.TrimSpace(out.String())
		if s.Out.Is != nil {
			assert.Equal(*s.Out.Is, outContents)
		}
		if len(s.Out.Contains) > 0 {
			for _, find := range s.Out.Contains {
				assert.Contains(outContents, find)
			}
		}
		if len(s.Out.ContainsOneOf) > 0 {
			found := false
			for _, find := range s.Out.ContainsOneOf {
				if idx := strings.Index(outContents, find); idx > -1 {
					found = true
					break
				}
			}
			if !found {
				assert.Fail(
					"expected to find one of %s in stdout.",
					s.Out.ContainsOneOf,
				)
			}
		}
	}
}
