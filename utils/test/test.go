package test

import (
	"errors"

	"github.com/rudderlabs/rudder-transformations/plugins"
	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

const TEST_PLUGIN_NAME = "test"

type TestPlugin struct{}

func (t *TestPlugin) doDefaultTransform(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["default"] = true
	return dataMap, nil
}

func (t *TestPlugin) doResponseTransform(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["response"] = true
	return dataMap, nil
}

func (t *TestPlugin) GetTransformer(data interface{}) (types.Transformer, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}

	switch dataMap["type"] {
	case "response":
		return types.TransformerFunc(t.doResponseTransform), nil
	default:
		return types.TransformerFunc(t.doDefaultTransform), nil
	}
}

func (t *TestPlugin) GetName() string {
	return TEST_PLUGIN_NAME
}

func GetTestPluginManager() *plugins.Manager {
	pluginManager := plugins.NewManager()
	pluginManager.AddPluginProvider(types.NewPluginProvider(destinations.PROVIDER_NAME))
	_ = pluginManager.AddPlugins(destinations.PROVIDER_NAME, &TestPlugin{}, destinations.DefaultPlugin)
	return pluginManager
}

var BadPlugin = types.NewSimplePlugin(
	"bad",
	func(data interface{}) (interface{}, error) {
		return nil, errors.New("bad plugin")
	},
)
