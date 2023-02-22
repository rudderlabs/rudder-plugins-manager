package test

import (
	"errors"

	"github.com/rudderlabs/rudder-transformations/plugins"
	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

const TEST_PLUGIN_NAME = "test"

type TestPlugin struct{}

func (p *TestPlugin) Name() string {
	return TEST_PLUGIN_NAME
}

func (p *TestPlugin) GetTransformer(data interface{}) (types.Transformer, error) {
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

func GetTestPluginManager() *plugins.Manager {
	pluginManager := plugins.NewManager()
	pluginManager.AddPluginProvider(destinations.NewPluginProvider())
	pluginManager.AddPlugins(destinations.PROVIDER_NAME, &TestPlugin{}, &destinations.DefaultPlugin{})
	return pluginManager
}
