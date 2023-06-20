// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"

	gdttypes "github.com/jaypipes/gdt-core/types"
)

type gdtContextKey string

var (
	pathKey     = gdtContextKey("gdt.path")
	pluginsKey  = gdtContextKey("gdt.plugins")
	fixturesKey = gdtContextKey("gdt.fixtures")
)

// ContextModifier sets some value on the context
type ContextModifier func(context.Context) context.Context

// WithPath sets a context's Path attribute
func WithPath(path string) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pathKey, path)
	}
}

// WithPlugins sets a context's Plugins attribute
func WithPlugins(plugins []gdttypes.Plugin) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pluginsKey, plugins)
	}
}

// WithFixtures sets a context's Fixtures attribute
func WithFixtures(fixtures []gdttypes.Fixture) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, fixturesKey, fixtures)
	}
}

// New returns a new Context
func New(mods ...ContextModifier) context.Context {
	ctx := context.TODO()
	for _, mod := range mods {
		ctx = mod(ctx)
	}
	return ctx
}
