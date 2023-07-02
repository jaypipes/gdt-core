// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	gjs "github.com/xeipuuv/gojsonschema"

	gdterrors "github.com/jaypipes/gdt-core/errors"
	gdttypes "github.com/jaypipes/gdt-core/types"
)

// Expect represents one or more assertions about JSON data responses
type Expect struct {
	// Length of the expected JSON string.
	Len *int `yaml:"len,omitempty"`
	// Paths is a map, keyed by JSONPath expression, of expected values to find
	// at that expression.
	Paths map[string]string `yaml:"paths,omitempty"`
	// PathFormats is a map, keyed by JSONPath expression, of expected formats
	// that values found at the expression should have.
	PathFormats map[string]string `yaml:"path_formats,omitempty"`
	// Schema is a file path to the JSONSchema that the JSON should validate
	// against.
	Schema string `yaml:"schema,omitempty"`
}

// New returns a `gdttypes.Assertions` that asserts various conditions about
// JSON content
func New(
	exp *Expect,
	content []byte,
) gdttypes.Assertions {
	return &assertions{
		failures: []error{},
		exp:      exp,
		content:  content,
	}
}

// assertions represents one or more assertions about JSON data responses and
// implements the gdttypes.Assertions interface
type assertions struct {
	// failures contains the set of error messages for failed assertions
	failures []error
	// terminal indicates there was a failure in evaluating the assertions that
	// should be considered a terminal condition (and therefore the test action
	// should not be retried).
	terminal bool
	// exp contains the expected conditions for to be asserted
	exp *Expect
	// content is the JSON content we will check
	content []byte
}

// Fail appends a supplied error to the set of failed assertions
func (a *assertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of failure messages indicating which assertions did
// not succeed.
func (a *assertions) Failures() []error {
	return a.failures
}

// Terminal returns true if re-executing the assertions against the same result
// would be pointless. This allows assertions to inform the Spec that retrying
// the same operation would not be necessary.
func (a *assertions) Terminal() bool {
	return a.terminal
}

// OK returns true if all contained assertions pass successfully
func (a *assertions) OK() bool {
	if a == nil || a.exp == nil {
		return true
	}
	if !a.lenOK() {
		return false
	}
	if !a.pathsOK() {
		return false
	}
	if !a.pathFormatsOK() {
		return false
	}
	if !a.schemaOK() {
		return false
	}
	return true
}

// lenOK returns true if the content length matches expectations, false
// otherwise
func (a *assertions) lenOK() bool {
	if a == nil || a.exp == nil {
		return true
	}
	if a.exp.Len != nil {
		exp := *a.exp.Len
		got := len(a.content)
		if exp != got {
			a.Fail(gdterrors.NotEqualLength(exp, got))
			return false
		}
	}
	return true
}

// pathsOK returns true if the content matches the Paths conditions, false
// otherwise
func (a *assertions) pathsOK() bool {
	if a == nil || a.exp == nil {
		return true
	}
	if len(a.exp.Paths) == 0 {
		return true
	}
	v := interface{}(nil)
	if err := json.Unmarshal(a.content, &v); err != nil {
		a.Fail(JSONUnmarshalError(err))
		a.terminal = true
		return false
	}
	for path, expVal := range a.exp.Paths {
		got, err := jsonpath.Get(path, v)
		if err != nil {
			a.Fail(JSONPathError(path, err))
			a.terminal = true
			return false
		}
		switch got.(type) {
		case string:
			if expVal != got.(string) {
				a.Fail(JSONPathNotEqual(path, expVal, got))
				return false
			}
		case int, uint, int64, uint64:
			expValInt, err := strconv.Atoi(expVal)
			if err != nil {
				a.Fail(JSONPathConversionError(path, expVal, got))
				a.terminal = true
				return false
			}
			if expValInt != got.(int) {
				a.Fail(JSONPathNotEqual(path, expVal, got))
				return false
			}
		case float32, float64:
			expValFloat, err := strconv.ParseFloat(expVal, 64)
			if err != nil {
				a.Fail(JSONPathConversionError(path, expVal, got))
				a.terminal = true
				return false
			}
			if expValFloat != got.(float64) {
				a.Fail(JSONPathNotEqual(path, expVal, got))
				return false
			}
		case bool:
			expValBool, err := strconv.ParseBool(expVal)
			if err != nil {
				a.Fail(JSONPathConversionError(path, expVal, got))
				a.terminal = true
				return false
			}
			if expValBool != got.(bool) {
				a.Fail(JSONPathNotEqual(path, expVal, got))
				return false
			}
		default:
			a.Fail(JSONPathConversionError(path, expVal, got))
			a.terminal = true
			return false
		}
	}
	return true
}

// pathFormatsOK returns true if the content matches the PathFormats
// conditions, false otherwise
func (a *assertions) pathFormatsOK() bool {
	if a == nil || a.exp == nil {
		return true
	}
	if len(a.exp.PathFormats) == 0 {
		return true
	}
	v := interface{}(nil)
	if e := json.Unmarshal(a.content, &v); e != nil {
		a.Fail(JSONUnmarshalError(e))
		a.terminal = true
		return false
	}
	for path, format := range a.exp.PathFormats {
		got, err := jsonpath.Get(path, v)
		if err != nil {
			a.Fail(JSONPathError(path, err))
			a.terminal = true
			return false
		}
		ok, err := isFormatted(format, got)
		if err != nil {
			a.Fail(JSONFormatError(format, err))
			a.terminal = true
			return false
		}
		if !ok {
			a.Fail(JSONFormatNotEqual(path, format))
			return false
		}
	}
	return true
}

// schemaOK returns true if the content matches the Schema condition, false
// otherwise
func (a *assertions) schemaOK() bool {
	if a == nil || a.exp == nil {
		return true
	}
	if a.exp.Schema == "" {
		return true
	}

	schemaPath := a.exp.Schema
	schemaLoader := gjs.NewReferenceLoader(schemaPath)
	docLoader := gjs.NewStringLoader(string(a.content))

	res, err := gjs.Validate(schemaLoader, docLoader)
	if err != nil {
		a.Fail(JSONSchemaValidateError(schemaPath, err))
		a.terminal = true
		return false
	}

	var errStr string
	if len(res.Errors()) > 0 {
		errStrs := make([]string, len(res.Errors()))
		for x, err := range res.Errors() {
			errStrs[x] = err.String()
		}
		errStr = "- " + strings.Join(errStrs, "\n- ")
	}
	if !res.Valid() {
		a.Fail(JSONSchemaInvalid(schemaPath, fmt.Errorf(errStr)))
	}
	return res.Valid()
}
