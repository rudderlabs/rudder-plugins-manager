package plugins_test

import (
	"testing"

	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestPluginManagerAddPlugin(t *testing.T) {
	pluginManager := test.GetTestPluginManager()
	plugin, err := pluginManager.GetPlugin(destinations.PROVIDER_NAME, test.TEST_PLUGIN_NAME)
	assert.Nil(t, err)
	assert.Equal(t, test.TEST_PLUGIN_NAME, plugin.GetName())

	data := map[string]any{
		"test": "test",
	}
	respTransformer, err := plugin.GetTransformer(map[string]any{
		"type": "response",
	})
	assert.Nil(t, err)
	transformedData, err := respTransformer.Transform(data)
	assert.Nil(t, err)
	assert.NotNil(t, transformedData.(map[string]any)["response"])

	defaultTransformer, err := plugin.GetTransformer(map[string]any{
		"type": "default",
	})
	assert.Nil(t, err)
	transformedData, err = defaultTransformer.Transform(data)
	assert.Nil(t, err)
	assert.NotNil(t, transformedData.(map[string]any)["default"])
}

func TestPluginManagerGetPlugin(t *testing.T) {
	pluginManager := test.GetTestPluginManager()
	plugin, err := pluginManager.GetPlugin("non-existing-provider", test.TEST_PLUGIN_NAME)
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "provider not found")

	plugin, err = pluginManager.GetPlugin(destinations.PROVIDER_NAME, "non-existing-plugin")
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "plugin not found")
}
