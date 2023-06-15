// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	// ErrUnknownSpec indicates that there was a test spec definition in a YAML
	// file that no plugin could parse.
	ErrUnknownSpec = errors.New("no plugin could parse spec definition")
	// ErrUnknownField indicates that there was an unknown field in the parsing
	// of a spec or scenario.
	ErrUnknownField = errors.New("unknown field")
	// ErrInvalid indicates a YAML definition is not valid
	ErrInvalid = errors.New("invalid YAML")
	// ErrInvalidExpectedMap indicates that we did not find an
	// expected mapping field
	ErrInvalidExpectedMap = fmt.Errorf(
		"%w: expected map field", ErrInvalid,
	)
	// ErrInvalidExpectedScalar indicates that we did not find an
	// expected scalar field
	ErrInvalidExpectedScalar = fmt.Errorf(
		"%w: expected scalar field", ErrInvalid,
	)
	// ErrInvalidExpectedSequence indicates that we did not find an
	// expected scalar field
	ErrInvalidExpectedSequence = fmt.Errorf(
		"%w: expected sequence field", ErrInvalid,
	)
	// ErrInvalidExpectedInt indicates that we did not find an
	// expected integer value
	ErrInvalidExpectedInt = fmt.Errorf(
		"%w: expected int value", ErrInvalid,
	)
)

// UnknownSpecAt returns an ErrUnknownSpec with the line/column of the supplied
// YAML node.
func UnknownSpecAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrUnknownSpec, node.Line, node.Column,
	)
}

// UnknownFieldAt returns an ErrUnknownField for a supplied field annotated
// with the line/column of the supplied YAML node.
func UnknownFieldAt(field string, node *yaml.Node) error {
	return fmt.Errorf(
		"%w: %q at line %d, column %d",
		ErrUnknownField, field, node.Line, node.Column,
	)
}

// ExpectedMapAt returns an ErrInvalidExpectedMap error annotated with the
// line/column of the supplied YAML node.
func ExpectedMapAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrInvalidExpectedMap, node.Line, node.Column,
	)
}

// ExpectedScalarAt returns an ErrInvalidExpectedScalar error annotated with
// the line/column of the supplied YAML node.
func ExpectedScalarAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrInvalidExpectedScalar, node.Line, node.Column,
	)
}

// ExpectedSequenceAt returns an ErrInvalidExpectedSequence error annotated
// with the line/column of the supplied YAML node.
func ExpectedSequenceAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrInvalidExpectedSequence, node.Line, node.Column,
	)
}

// ExpectedIntAt returns an ErrInvalidExpectedInt error annotated
// with the line/column of the supplied YAML node.
func ExpectedIntAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrInvalidExpectedInt, node.Line, node.Column,
	)
}
