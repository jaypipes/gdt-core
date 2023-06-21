// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"errors"

	"github.com/jaypipes/gdt-core/context"
	gdterrors "github.com/jaypipes/gdt-core/errors"
	gdtexec "github.com/jaypipes/gdt-core/exec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"gopkg.in/yaml.v3"
)

// coreSpecPrototypes returns a slice of known non-plugin Spec types
func coreSpecPrototypes() []gdttypes.Spec {
	return []gdttypes.Spec{
		&gdtexec.ExecSpec{},
	}
}

// UnmarshalYAML is a custom unmarshaler that asks plugins for their known spec
// types and attempts to unmarshal test spec contents into those types.
func (s *Scenario) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return gdterrors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return gdterrors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return gdterrors.ExpectedScalarAt(valNode)
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return gdterrors.ExpectedScalarAt(valNode)
			}
			s.Description = valNode.Value
		case "require":
			if valNode.Kind != yaml.SequenceNode {
				return gdterrors.ExpectedSequenceAt(valNode)
			}
			requires := make([]string, len(valNode.Content))
			for x, n := range valNode.Content {
				requires[x] = n.Value
			}
			s.Require = requires
		case "defaults":
			if valNode.Kind != yaml.MappingNode {
				return gdterrors.ExpectedMapAt(valNode)
			}
			defaults := map[string]interface{}{}
			if err := valNode.Decode(&defaults); err != nil {
				return err
			}
			s.Defaults = defaults
		case "tests":
			if valNode.Kind != yaml.SequenceNode {
				return gdterrors.ExpectedSequenceAt(valNode)
			}
			for idx, testNode := range valNode.Content {
				parsed := false
				specs := coreSpecPrototypes()
				for _, p := range context.Plugins(s.ctx) {
					specs = append(specs, p.Specs()...)
				}
				for _, sp := range specs {
					if err := testNode.Decode(sp); err != nil {
						if errors.Is(err, gdterrors.ErrUnknownField) {
							continue
						}
						return err
					} else {
						if err := sp.SetBaseFields(idx, testNode); err != nil {
							return err
						}
						s.Tests = append(s.Tests, sp)
						parsed = true
						break
					}
				}
				if !parsed {
					return gdterrors.UnknownSpecAt(valNode)
				}
			}
		}

	}
	return nil
}
