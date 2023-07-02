// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	gdtjson "github.com/jaypipes/gdt-core/assertion/json"
	gdterrors "github.com/jaypipes/gdt-core/errors"
	"github.com/stretchr/testify/require"
)

func content() []byte {
	b, _ := ioutil.ReadFile(filepath.Join("testdata", "books.json"))
	return b
}

func TestLength(t *testing.T) {
	require := require.New(t)

	c := content()
	expLen := len(c)

	exp := gdtjson.Expect{
		Len: &expLen,
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK())
	require.False(a.Terminal())
	require.Empty(a.Failures())

	expLen = 0
	a = gdtjson.New(&exp, c)
	require.False(a.OK())
	require.False(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdterrors.ErrNotEqual)
}

func TestJSONUnmarshalError(t *testing.T) {
	require := require.New(t)

	c := []byte(`not { value } json`)

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"1234": "foo",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK())
	require.True(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONUnmarshalError)
}

func TestJSONPathError(t *testing.T) {
	require := require.New(t)

	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			// This is not a valid JSONPath expression... must begin with the
			// root element $
			"[0].pages": "127",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK())
	require.True(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathError)
}

func TestJSONPathConversionError(t *testing.T) {
	require := require.New(t)

	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"1234": "foo",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK())
	require.True(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathConversionError)
}

func TestJSONPathNotEqual(t *testing.T) {
	require := require.New(t)

	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"$[0].pages": "127",
		},
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK())
	require.False(a.Terminal())
	require.Empty(a.Failures())

	exp = gdtjson.Expect{
		Paths: map[string]string{
			"$[0].pages": "42",
		},
	}

	a = gdtjson.New(&exp, c)
	require.False(a.OK())
	require.False(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathNotEqual)
}

func TestJSONPathFormatError(t *testing.T) {
	require := require.New(t)

	c := content()

	exp := gdtjson.Expect{
		PathFormats: map[string]string{
			"$[0].pages": "invalidformat",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK())
	require.True(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONFormatError)
}

func TestJSONPathFormatNotEqual(t *testing.T) {
	require := require.New(t)

	c := content()

	exp := gdtjson.Expect{
		PathFormats: map[string]string{
			"$[0].id": "uuid4",
		},
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK())
	require.False(a.Terminal())
	require.Empty(a.Failures())

	exp = gdtjson.Expect{
		PathFormats: map[string]string{
			"$[0].pages": "uuid4",
		},
	}

	a = gdtjson.New(&exp, c)
	require.False(a.OK())
	require.False(a.Terminal())
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONFormatNotEqual)
}
