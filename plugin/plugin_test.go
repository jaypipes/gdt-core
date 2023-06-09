// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin_test

import (
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/stretchr/testify/assert"
)

type fooPlugin struct{}

func (p *fooPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "foo",
	}
}

func (p *fooPlugin) Parse(gdttypes.Appendable, []byte) error {
	return nil
}

func TestRegisterAndList(t *testing.T) {
	assert := assert.New(t)

	plugins := plugin.List()
	assert.Equal(0, len(plugins))

	plugin.Register(&fooPlugin{})

	plugins = plugin.List()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)

	// Register called twice with the same named plugin should be be a no-op

	plugin.Register(&fooPlugin{})

	plugins = plugin.List()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)
}
