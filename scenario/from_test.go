// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/gdt-core/plugin"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/spec"
	gdttypes "github.com/jaypipes/gdt-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromNoPlugins(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "http-failure.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	sc := s.(*scenario.Scenario)
	assert.Equal("HTTP failure", sc.Name)
	assert.Equal("testdata/http-failure.yaml", sc.Path)
	assert.Equal([]string{"books_api", "books_data"}, sc.Require)
	assert.Equal(
		map[string]interface{}{
			"http": map[string]interface{}{
				"base_url": "http://127.0.0.1:4000",
			},
		},
		sc.Defaults,
	)
	assert.Equal(
		[]*spec.Spec{
			&spec.Spec{
				Name:        "no such book was found",
				Description: "",
			},
		},
		sc.Tests,
	)
}

type failingPlugin struct{}

func (p *failingPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name: "failer",
	}
}

func (p *failingPlugin) Parse(gdttypes.Appendable, []byte) error {
	return fmt.Errorf("Indy, bad dates!")
}

func TestFromFailingPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "http-failure.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	plugin.Register(&failingPlugin{})

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.Nil(s)
}
