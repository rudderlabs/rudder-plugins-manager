package plugins_test

import (
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-transformations/plugins"
	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/plugins/types"
	"github.com/stretchr/testify/assert"
)

func TestPluginManager(t *testing.T) {
	pluginManager := plugins.ManagerInstance
	plugin, ok := pluginManager.GetPlugin(destinations.PROVIDER_NAME, "default")
	assert.True(t, ok)
	assert.Equal(t, "default", plugin.Name())
	data := map[string]interface{}{
		"test": "test",
	}
	transformer, err := plugin.GetTransformer(data)
	assert.Nil(t, err)
	transformedData, err := transformer.Transform(data)
	assert.Nil(t, err)
	assert.Equal(t, data, transformedData)
}

type testPlugin struct{}

func (p *testPlugin) Name() string {
	return "test"
}
func (p *testPlugin) GetTransformer(data interface{}) (types.Transformer, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}

	switch dataMap["type"] {
	case "response":
		return types.TransformerFunc(func(data interface{}) (interface{}, error) {
			data.(map[string]interface{})["response"] = true
			return data, nil
		}), nil
	default:
		return types.TransformerFunc(func(data interface{}) (interface{}, error) {
			data.(map[string]interface{})["default"] = true
			return data, nil
		}), nil
	}
}

func TestPluginManagerAddPlugin(t *testing.T) {
	pluginManager := plugins.ManagerInstance
	testPlugin := &testPlugin{}
	pluginManager.AddPlugin(destinations.PROVIDER_NAME, testPlugin)
	plugin, ok := pluginManager.GetPlugin(destinations.PROVIDER_NAME, testPlugin.Name())
	assert.True(t, ok)
	assert.Equal(t, testPlugin.Name(), plugin.Name())

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
