// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"
	"io"

	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/samber/lo"
)

type ContextKey string

var (
	debugKey    = ContextKey("gdt.debug")
	pluginsKey  = ContextKey("gdt.plugins")
	fixturesKey = ContextKey("gdt.fixtures")
	priorRunKey = ContextKey("gdt.run.prior")
)

// ContextModifier sets some value on the context
type ContextModifier func(context.Context) context.Context

// WithDebug sets a context's Debug writer. If you want gdt to log extra
// debugging information about tests and assertions, pass it a context with a
// debug writer:
//
// ```go
// f := ioutil.TempFile("", "mytest*.log")
// ctx := gdtcontext.New(gdtcontext.WithDebug(f))
// ```
//
// you can then inspect the debug "log" and do whatever you'd like with it.
//
// Or you could pass a console writer and just have gdt write to the console
// its debugging information:
//
// ```go
// ctx := gdtcontext.New(gdtcontext.WithDebug(os.Stdout))
// ```
func WithDebug(debug io.Writer) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, debugKey, debug)
	}
}

// WithPlugins sets a context's Plugins
func WithPlugins(plugins []gdttypes.Plugin) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pluginsKey, plugins)
	}
}

// WithFixtures sets a context's Fixtures
func WithFixtures(fixtures map[string]gdttypes.Fixture) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, fixturesKey, fixtures)
	}
}

// SetDebug sets gdt's debug logging to the supplied `io.Writer`
func SetDebug(
	ctx context.Context,
	debug io.Writer,
) context.Context {
	return context.WithValue(ctx, debugKey, debug)
}

// RegisterFixture registers a named fixtures with the context
func RegisterFixture(
	ctx context.Context,
	name string,
	f gdttypes.Fixture,
) context.Context {
	fixtures := Fixtures(ctx)
	fixtures[name] = f
	return context.WithValue(ctx, fixturesKey, fixtures)
}

// RegisterPlugin registers a plugin with the context
func RegisterPlugin(
	ctx context.Context,
	p gdttypes.Plugin,
) context.Context {
	plugins := Plugins(ctx)
	for _, plug := range plugins {
		if plug.Info().Name == p.Info().Name {
			// No need to register... already known.
			return ctx
		}
	}
	plugins = append(plugins, p)
	return context.WithValue(ctx, pluginsKey, plugins)
}

// StorePriorRun saves prior run data in the context. If there is already prior
// run data cached in the supplied context, the existing data is merged with
// the supplied data.
func StorePriorRun(
	ctx context.Context,
	data map[string]interface{},
) context.Context {
	existing := PriorRun(ctx)
	merged := lo.Assign(existing, data)
	return context.WithValue(ctx, priorRunKey, merged)
}

// New returns a new Context
func New(mods ...ContextModifier) context.Context {
	ctx := context.TODO()
	for _, mod := range mods {
		ctx = mod(ctx)
	}
	return ctx
}
