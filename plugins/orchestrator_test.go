package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/stretchr/testify/assert"
)

var orchestrator = plugins.NewTransformPlugin("orchestrator", func(data *plugins.Message) (*plugins.Message, error) {
	dataMap, ok := data.Data.(map[string]any)
	if !ok {
		return nil, errors.New("data is not a map")
	}
	if dataMap["test"] != nil {
		return plugins.NextPluginMessage("test"), nil
	} else if dataMap["secret"] != nil {
		return plugins.NextPluginMessage("blobl"), nil
	}
	return nil, nil
})

func TestOrchestratorPlugin(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	identity := plugins.NewTransformPlugin("identity", func(data *plugins.Message) (*plugins.Message, error) {
		return data, nil
	})
	manager.Add(identity)
	orchestrator := plugins.NewTransformPlugin("orchestrator",
		func(data *plugins.Message) (*plugins.Message, error) {
			return plugins.NextPluginMessage("identity"), nil
		},
	)
	pluginOrchestrator := plugins.NewOrchestratorPlugin(manager, orchestrator)
	assert.NotNil(t, pluginOrchestrator)
	data, err := pluginOrchestrator.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, plugins.NewMessage(map[string]any{}), data)
}

func TestOrchestratorPluginFailureCases(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	orchestrator := plugins.NewTransformPlugin("orchestrator", func(data *plugins.Message) (*plugins.Message, error) {
		dataMap, ok := data.Data.(map[string]any)
		if !ok {
			return nil, errors.New("data is not a map")
		}
		if dataMap["test"] != nil {
			return plugins.NextPluginMessage("test"), nil
		} else {
			return nil, errors.New("invalid input")
		}
	})

	manager.AddOrchestrator(orchestrator)
	pluginOrchestrator, err := manager.Get(orchestrator.GetName())
	assert.Nil(t, err)
	data, err := pluginOrchestrator.Execute(context.Background(), testMessage())
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "plugin not found")

	data, err = pluginOrchestrator.Execute(context.Background(), emptyMessage())
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "invalid input")

	manager.AddOrchestrator(testPlugin)
	pluginOrchestrator, err = manager.Get(testPlugin.GetName())
	assert.Nil(t, err)
	data, err = pluginOrchestrator.Execute(context.Background(), testMessage())
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "failed to get next plugin name")
}
