package gdtcore

import "gopkg.in/yaml.v3"

// Spec represents things that can be parsed from YAML nodes and can have
// their Name and Description fields set.
type Spec interface {
	Runnable
	yaml.Unmarshaler
	SetBaseFields(*yaml.Node) error
}
