package plugins_test

import (
	"testing"

	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestPluginManagerAddPlugin(t *testing.T) {
	pluginManager := test.GetTestPluginManager()
	plugin, ok := pluginManager.GetPlugin(destinations.PROVIDER_NAME, test.TEST_PLUGIN_NAME)
	assert.True(t, ok)
	assert.Equal(t, test.TEST_PLUGIN_NAME, plugin.Name())

	data := map[string]interface{}{
		"test": "test",
	}
	respTransformer, err := plugin.GetTransformer(map[string]interface{}{
		"type": "response",
	})
	assert.Nil(t, err)
	transformedData, err := respTransformer.Transform(data)
	assert.Nil(t, err)
	assert.NotNil(t, transformedData.(map[string]interface{})["response"])

	defaultTransformer, err := plugin.GetTransformer(map[string]interface{}{
		"type": "default",
	})
	assert.Nil(t, err)
	transformedData, err = defaultTransformer.Transform(data)
	assert.Nil(t, err)
	assert.NotNil(t, transformedData.(map[string]interface{})["default"])
}
