package plugins_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var testWorkflow = lo.Must(plugins.NewBaseWorkflowPlugin(nil, plugins.WorkflowConfig{
	Name: "test",
	Steps: []plugins.StepConfig{
		{
			Name:     "test",
			Bloblang: `root.test = "test"`,
		},
	},
}))

func TestWorkflowManager(t *testing.T) {
	manager := plugins.NewBaseWorkflowManager()
	manager.Add(testWorkflow)

	workflow, err := manager.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, testWorkflow, workflow)
	assert.True(t, manager.Has("test"))

	_, err = manager.Get("non-existent")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "workflow not found")

	assert.False(t, manager.Has("non-existent"))

	data, err := manager.Execute(context.Background(), "test", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, data.Data)
}

func TestNewBasePluginManager(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	manager.Add(testPlugin)
	plugin, err := manager.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", plugin.GetName())
	assert.True(t, manager.Has("test"))

	plugin, err = manager.Get("non-existing-plugin")
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "plugin not found")
	assert.False(t, manager.Has("non-existing-plugin"))

	result, err := manager.Execute(context.Background(), "test", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), result)

	result, err = manager.Execute(context.Background(), "non-existing-plugin", emptyMessage())
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "plugin not found")
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
