// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package spec

import (
	gdttypes "github.com/jaypipes/gdt-core/types"
	"gopkg.in/yaml.v3"
)

// FromBytes returns a Spec after parsing the supplied contents
func FromBytes(contents []byte) (gdttypes.Runnable, error) {
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" top-level element in the
	// YAML. If no "type" element was found, the test type defaults to HTTP.
	// Once the type is determined, then the test case module (e.g. gdt/http)
	// is called to parse the file into the case type-specific schema
	s := New()
	if err := yaml.Unmarshal(contents, s); err != nil {
		return nil, err
	}

	return s, nil
}
