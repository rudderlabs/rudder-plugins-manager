package types_test

import (
	"testing"

	"github.com/rudderlabs/rudder-transformations/plugins/types"
	"github.com/rudderlabs/rudder-transformations/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestNewPluginProvider(t *testing.T) {
	pluginProvider := types.NewPluginProvider("test")
	assert.Equal(t, "test", pluginProvider.Name)
	assert.Equal(t, "test", pluginProvider.GetName())
	pluginProvider.AddPlugin(&test.TestPlugin{})
	plugin, err := pluginProvider.GetPlugin("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", plugin.GetName())
	plugin, err = pluginProvider.GetPlugin("non-existing-plugin")
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "plugin not found")
}
