// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jaypipes/gdt-core/errors"
)

// errUnknownShell returns a wrapped version of ErrInvalid that indicates the
// user specified an unknown shell.
func errUnknownShell(shell string) error {
	return fmt.Errorf(
		"%w: expected map field", errors.ErrInvalid,
	)
}

func (s *ExecSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "shell":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Shell = strings.TrimSpace(valNode.Value)
			if _, err := exec.LookPath(s.Shell); err != nil {
				return errUnknownShell(s.Shell)
			}
		case "exec":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Exec = strings.TrimSpace(valNode.Value)
		case "exit_code":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			ec, err := strconv.Atoi(valNode.Value)
			if err != nil {
				return err
			}
			s.ExitCode = ec
		case "out":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			var out *PipeAssertions
			if err := valNode.Decode(&out); err != nil {
				return err
			}
			s.Out = out
		case "name", "description":
			continue
		default:
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}
