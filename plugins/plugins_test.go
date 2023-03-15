package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var testPlugin = plugins.NewTransformPlugin("test", func(msg *plugins.Message) (*plugins.Message, error) {
	dataMap, ok := msg.Data.(map[string]any)
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["test"] = "test"
	return plugins.NewMessage(dataMap), nil
})

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

func emptyMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{})
}

func testMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"test": "test"})
}

func complexMessage() *plugins.Message {
	type linkedList struct {
		Data int
		Next *linkedList
	}
	return plugins.NewMessage(linkedList{
		Data: 1,
		Next: &linkedList{
			Data: 1,
		},
	})
}

func secretMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"secret": "secret"})
}

func secretEncodedMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"secret": "c2VjcmV0"})
}

var badPlugin = plugins.NewTransformPlugin("bad", func(msg *plugins.Message) (*plugins.Message, error) {
	return nil, errors.New("bad plugin")
})

var bloblPlugin = lo.Must(plugins.NewBloblangPlugin("blobl", `root.secret = this.data.secret.encode("base64")`))

func TestNewSimplePlugin(t *testing.T) {
	data, err := testPlugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestNewBasePlugin(t *testing.T) {
	plugin := plugins.NewBasePlugin("test",
		plugins.ExecuteFunc(
			func(ctx context.Context, data *plugins.Message) (*plugins.Message, error) {
				return data, nil
			},
		),
	)
	assert.NotNil(t, plugin)
	data, err := plugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, emptyMessage(), data)
	assert.Equal(t, "test", plugin.GetName())
}

func TestNewSimpleBlobLPlugin(t *testing.T) {
	data, err := bloblPlugin.Execute(context.Background(), secretMessage())

	assert.Nil(t, err)
	assert.Equal(t, secretEncodedMessage(), data)

	plugin, err := plugins.NewBloblangPlugin("test", `
	root.data = this.data
	root.data.test = "test"
	`)
	assert.Nil(t, err)
	data, err = plugin.Execute(context.Background(), complexMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{
		"Data": float64(1),
		"Next": map[string]any{
			"Data": float64(1),
			"Next": nil,
		},
		"test": "test",
	}, data.Data)
}

func TestNewSimpleBlobLPluginFailureCase(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `some invalid blobl`)
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "failed to parse bloblang template")

	plugin, err = plugins.NewBloblangPlugin("test", `throw("some error")`)
	assert.Nil(t, err)
	_, err = plugin.Execute(context.Background(), emptyMessage())
	assert.ErrorContains(t, err, "some error")
}

func TestNewSimpleBlobLPluginCondition(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `this.data.test == "test"`)
	assert.Nil(t, err)
	data, err := plugin.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.True(t, data.Data.(bool))
}

func TestNewBasePluginManager(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	manager.Add(testPlugin)
	plugin, err := manager.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", plugin.GetName())
	plugin, err = manager.Get("non-existing-plugin")
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "plugin not found")

	result, err := manager.Execute(context.Background(), "test", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), result)

	result, err = manager.Execute(context.Background(), "non-existing-plugin", emptyMessage())
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "plugin not found")
}

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

func TestPluginManagerAddOrchestrator(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	manager.Add(testPlugin)
	manager.Add(bloblPlugin)

	manager.AddOrchestrator(orchestrator)
	pluginOrchestrator, err := manager.Get("orchestrator")
	assert.Nil(t, err)
	data, err := pluginOrchestrator.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
	data, err = pluginOrchestrator.Execute(context.Background(), secretMessage())
	assert.Nil(t, err)
	assert.Equal(t, secretEncodedMessage(), data)
	data, err = pluginOrchestrator.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Nil(t, data)
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
