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

func TestPluginManager(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	manager.Add(testPlugin)
	manager.AddOrchestrator(orchestrator)

	plugin, err := manager.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, testPlugin, plugin)

	data, err := manager.Execute(context.Background(), "test", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestWorkflowManager(t *testing.T) {
	manager := plugins.NewBaseWorkflowManager()
	manager.Add(testWorkflow)

	workflow, err := manager.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, testWorkflow, workflow)

	data, err := manager.Execute(context.Background(), "test", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test":"test"}, data.Data)
}
