// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"io"
	"io/ioutil"

	"github.com/jaypipes/gdt-core/plugin"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"gopkg.in/yaml.v3"
)

// FromReader parses the supplied io.Reader and returns a Scenario representing
// the contents in the reader. Returns an error if any syntax or validation
// failed
func FromReader(
	r io.Reader,
	mods ...ScenarioModifier,
) (gdttypes.Runnable, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(contents, mods...)
}

// FromBytes returns a Scenario after parsing the supplied contents
func FromBytes(
	contents []byte,
	mods ...ScenarioModifier,
) (gdttypes.Runnable, error) {
	// We do a double-parse of the test scenario file. The first pass
	// determines the type of test by simply looking for a "type" top-level
	// element in the YAML. If no "type" element was found, the test type
	// defaults to HTTP.  Once the type is determined, then the test case
	// module (e.g. gdt-http) is called to parse the file into the case
	// type-specific schema
	s := New(mods...)
	if err := yaml.Unmarshal(contents, s); err != nil {
		return nil, err
	}

	for _, p := range plugin.List() {
		if err := p.Parse(s, contents); err != nil {
			return nil, err
		}

	}

	return s, nil
}
