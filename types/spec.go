// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import "gopkg.in/yaml.v3"

// Spec represents things that can be parsed from YAML nodes and can have
// their Name and Description fields set.
type Spec interface {
	Runnable
	yaml.Unmarshaler
	// SetBaseFields sets the index for the Spec and examines the mapping YAML
	// node for a name and description field and sets the associated
	// Name/Description struct field from that value node.
	SetBaseFields(int, *yaml.Node) error
	// Title returns the human-readable name of the Spec. Default
	// implementation of this returns the Name of the Spec, if not nil, or a
	// slugified version of the Description, if present, falling back to the
	// Spec index if neither Name or Description are present.
	Title() string
}
