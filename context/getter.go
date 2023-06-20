// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"

	gdttypes "github.com/jaypipes/gdt-core/types"
)

// Path gets a context's Path
func Path(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(pathKey); v != nil {
		return v.(string)
	}
	return ""
}

// Plugins gets a context's Plugins
func Plugins(ctx context.Context) []gdttypes.Plugin {
	if ctx == nil {
		return []gdttypes.Plugin{}
	}
	if v := ctx.Value(pluginsKey); v != nil {
		return v.([]gdttypes.Plugin)
	}
	return []gdttypes.Plugin{}
}

// Fixtures gets a context's Fixtures
func Fixtures(ctx context.Context) map[string]gdttypes.Fixture {
	if ctx == nil {
		return map[string]gdttypes.Fixture{}
	}
	if v := ctx.Value(fixturesKey); v != nil {
		return v.(map[string]gdttypes.Fixture)
	}
	return map[string]gdttypes.Fixture{}
}
