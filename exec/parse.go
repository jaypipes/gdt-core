// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jaypipes/gdt-core/errors"
	gdttypes "github.com/jaypipes/gdt-core/types"
)

// errUnknownShell returns a wrapped version of ErrInvalid that indicates the
// user specified an unknown shell.
func errUnknownShell(shell string) error {
	return fmt.Errorf(
		"%w: expected map field", errors.ErrInvalid,
	)
}

func (s *Spec) UnmarshalYAML(node *yaml.Node) error {
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
			var pa *PipeAssertions
			if err := valNode.Decode(&pa); err != nil {
				return err
			}
			s.Out = pa
		case "err":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			var pa *PipeAssertions
			if err := valNode.Decode(&pa); err != nil {
				return err
			}
			s.Err = pa
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}
