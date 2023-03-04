package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/rudderlabs/rudder-plugins-manager/utils"
	"github.com/stretchr/testify/assert"
)

var testPlugin = plugins.NewTransformPlugin("test", func(data any) (any, error) {
	dataMap, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["test"] = "test"
	return dataMap, nil
})

var sampleStepPlugin = plugins.BaseStepPlugin{
	Name: "sampleStepPlugin",
	ConditionFunc: func(data plugins.StepInput) (bool, error) {
		return true, nil
	},
	ExecuteFunc: func(ctx context.Context, input plugins.StepInput) (*plugins.StepOutput, error) {
		return &plugins.StepOutput{
			Output:  map[string]interface{}{"test": "test"},
			Context: map[string]any{},
		}, nil
	},
}

var sampleBloblangStepPlugin = utils.Must(plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
	Name:      "sampleBloblangStepPlugin",
	Condition: utils.StringPtr(`this.data.test == "test"`),
	Template:  `root.output.blobl = true`,
}))

func TestNewSimplePlugin(t *testing.T) {
	data, err := testPlugin.Execute(context.Background(), map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, data)
}

func TestNewSimpleBlobLPlugin(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `root.secret = this.secret.encode("base64")`)
	assert.Nil(t, err)
	data, err := plugin.Execute(context.Background(), map[string]any{
		"secret": "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"secret": "c2VjcmV0"}, data)
}

func TestNewSimpleBlobLPluginCondition(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `this.test == "test"`)
	assert.Nil(t, err)
	data, err := plugin.Execute(context.Background(), map[string]any{
		"test": "test",
	})

	assert.Nil(t, err)
	assert.True(t, data.(bool))

	plugin, err = plugins.NewBloblangPlugin("test", `this.data.test == "test"`)
	assert.Nil(t, err)
	stepInput := plugins.StepInput{
		Data: map[string]any{
			"test": "test",
		},
	}
	data, err = plugin.Execute(context.Background(), stepInput.ToMap())

	assert.Nil(t, err)
	assert.True(t, data.(bool))
}

func TestNewPluginManager(t *testing.T) {
	manager := plugins.NewPluginManager()
	manager.AddPlugin(testPlugin)
	plugin, err := manager.GetPlugin("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", plugin.GetName())
	plugin, err = manager.GetPlugin("non-existing-plugin")
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "plugin not found")
}

func TestWorkflowPlugin(t *testing.T) {
	plugin := plugins.NewWorkflowPlugin("test", &sampleStepPlugin, sampleBloblangStepPlugin)
	assert.NotNil(t, plugin)
	input := map[string]any{}
	data, err := plugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"blobl": true}, data)
}
