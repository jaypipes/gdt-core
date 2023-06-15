// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec

import (
	"gopkg.in/yaml.v3"

	"github.com/jaypipes/gdt-core/errors"
)

// SetNameAndDescription examines the mapping YAML node for a name and
// description field and sets the associated Name/Description struct field from
// that value node.
func (s *Spec) SetNameAndDescription(node *yaml.Node) error {
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
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Description = valNode.Value
		}
	}
	return nil
}
