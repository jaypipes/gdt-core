// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	gdttypes "github.com/jaypipes/gdt-core/types"
	"gopkg.in/yaml.v3"
)

const (
	// hopefully nobody actually has an environment variable with this key!
	dollarSignReplacementToken = "oiuqdfjhaso7t213041"
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
	s := New(mods...)
	expanded := expandWithFixedDoubleDollar(string(contents))
	if err := yaml.Unmarshal([]byte(expanded), s); err != nil {
		return nil, err
	}

	return s, nil
}

// expandWithFixedDoubleDollar expands the given string using os.ExpandEnv,
// however unlike the default behaviour of replacing a string "$$VALUE" with
// "VALUE", it replaces the "$$" witha single "$". This allows test authors to
// use the dollar symbol in their test contents (they need to escape with
// '$$').
func expandWithFixedDoubleDollar(subject string) string {
	os.Setenv(dollarSignReplacementToken, "$")
	replaceStr := fmt.Sprintf("${%s}", dollarSignReplacementToken)
	return os.ExpandEnv(strings.Replace(subject, "$$", replaceStr, -1))
}
