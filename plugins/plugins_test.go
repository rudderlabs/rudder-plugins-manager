package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/rudderlabs/rudder-plugins-manager/types"
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

var bloblPlugin = utils.Must(plugins.NewBloblangPlugin("blobl", `root.secret = this.secret.encode("base64")`))

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

func TestNewBasePlugin(t *testing.T) {
	plugin := plugins.NewBasePlugin("test", types.ExecuteFunc(func(ctx context.Context, data any) (any, error) {
		return data, nil
	}))
	assert.NotNil(t, plugin)
	data, err := plugin.Execute(context.Background(), map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{}, data)
	assert.Equal(t, "test", plugin.GetName())
}

func TestNewSimpleBlobLPlugin(t *testing.T) {
	data, err := bloblPlugin.Execute(context.Background(), map[string]any{
		"secret": "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"secret": "c2VjcmV0"}, data)
}

func TestNewSimpleBlobLPluginFailureCase(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `some invalid blobl`)
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "failed to parse bloblang template")
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

	result, err := manager.Execute(context.Background(), "test", map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, result)

	result, err = manager.Execute(context.Background(), "non-existing-plugin", map[string]any{})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "plugin not found")
}

func TestStepPlugin(t *testing.T) {
	stepPlugin := plugins.BaseStepPlugin{Name: "stepPlugin"}
	assert.Equal(t, stepPlugin.Name, stepPlugin.GetName())
	shouldExecute, err := stepPlugin.ShouldExecute(plugins.StepInput{Data: map[string]any{}})
	assert.Nil(t, err)
	assert.Equal(t, true, shouldExecute)

	data, err := stepPlugin.Execute(context.Background(), plugins.StepInput{Data: map[string]any{}})
	assert.Nil(t, err)
	assert.Nil(t, data)
}

func TestBloblangStepPlugin(t *testing.T) {
	plugin, err := plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:     "bloblStepPlugin",
		Template: `root.output.blobl = true`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, plugin)
	assert.Equal(t, "bloblStepPlugin", plugin.GetName())
	shouldExecute, err := plugin.ShouldExecute(plugins.StepInput{Data: map[string]any{}})
	assert.Nil(t, err)
	assert.Equal(t, true, shouldExecute)
	data, err := plugin.Execute(context.Background(), plugins.StepInput{Data: map[string]any{}})
	assert.Nil(t, err)
	assert.Equal(t, &plugins.StepOutput{Output: map[string]any{"blobl": true}}, data)
}

func TestBloblangStepPluginFailureCases(t *testing.T) {
	_, err := plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:      "bloblStepPlugin",
		Condition: utils.StringPtr("bad condition"),
		Template:  `root.output.blobl = true`,
	})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to parse condition")

	_, err = plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:     "bloblStepPlugin",
		Template: `bad template`,
	})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to parse template")

	plugin, err := plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:      "bloblStepPlugin",
		Condition: utils.StringPtr(`throw("error")`),
		Template:  `throw("error")`,
	})
	assert.Nil(t, err)
	_, err = plugin.ShouldExecute(plugins.StepInput{Data: map[string]any{}})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "error")

	plugin, err = plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:     "bloblStepPlugin",
		Template: `throw("error")`,
	})
	assert.Nil(t, err)
	_, err = plugin.Execute(context.Background(), plugins.StepInput{Data: map[string]any{}})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "error")

	plugin, err = plugins.NewBloblangStepPlugin(plugins.BloblangWorkflowStep{
		Name:     "bloblStepPlugin",
		Template: `root = "hello world"`,
	})
	assert.Nil(t, err)
	_, err = plugin.Execute(context.Background(), plugins.StepInput{Data: map[string]any{}})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "expected a map")
}

func TestWorkflowPlugin(t *testing.T) {
	plugin := plugins.NewWorkflowPlugin("workflow", &sampleStepPlugin, sampleBloblangStepPlugin)
	assert.NotNil(t, plugin)
	assert.Equal(t, "workflow", plugin.GetName())
	input := map[string]any{}
	data, err := plugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"blobl": true}, data)
}

func TestBloblangWorkflowPlugin(t *testing.T) {
	plugin, err := plugins.NewBloblangWorkflowPlugin(&plugins.BloblangWorkflow{
		Name: "workflow",
		Steps: []plugins.BloblangWorkflowStep{
			{
				Name:      "step1",
				Condition: utils.StringPtr(`this.data.test == "test"`),
				Template:  `root.output.blobl = true`,
			},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, plugin)
	assert.Equal(t, "workflow", plugin.GetName())
	input := map[string]any{"test": "test"}
	data, err := plugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"blobl": true}, data)
}

func TestBloblangWorkflowPluginFailureCases(t *testing.T) {
	plugin, err := plugins.NewBloblangWorkflowPlugin(&plugins.BloblangWorkflow{
		Name: "workflow",
		Steps: []plugins.BloblangWorkflowStep{
			{
				Name:     "step1",
				Template: `bad template`,
			},
		},
	})
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "failed to parse template")
}

func TestWorkflowPluginFailureCases(t *testing.T) {
	badStepPlugin := plugins.BaseStepPlugin{Name: "badStepPlugin", ConditionFunc: func(data plugins.StepInput) (bool, error) {
		return false, errors.New("error")
	}}
	workflowPlugin := plugins.NewWorkflowPlugin("workflow", &badStepPlugin)
	input := map[string]any{}
	data, err := workflowPlugin.Execute(context.Background(), input)
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.EqualError(t, err, "error")
	badStepPlugin = plugins.BaseStepPlugin{Name: "badStepPlugin", ExecuteFunc: func(_ context.Context, data plugins.StepInput) (*plugins.StepOutput, error) {
		return nil, errors.New("error")
	}}
	data, err = workflowPlugin.Execute(context.Background(), input)
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.EqualError(t, err, "error")
}

func TestOrchestratorPlugin(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	identity := plugins.NewTransformPlugin("identity", func(data any) (any, error) {
		return data, nil
	})
	manager.Add(identity)
	orchestrator := plugins.NewTransformPlugin("orchestrator", func(data any) (any, error) {
		return "identity", nil
	})
	pluginOrchestrator := plugins.NewOrchestratorPlugin(manager, orchestrator)
	assert.NotNil(t, pluginOrchestrator)
	data, err := pluginOrchestrator.Execute(context.Background(), map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{}, data)
}

func TestPluginManagerAddOrchestrator(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	manager.Add(testPlugin)
	manager.Add(bloblPlugin)
	orchestrator := plugins.NewTransformPlugin("orchestrator", func(data any) (any, error) {
		dataMap, ok := data.(map[string]any)
		if !ok {
			return nil, errors.New("data is not a map")
		}
		if dataMap["testPlugin"] != nil {
			return "test", nil
		}
		return "blobl", nil
	})

	manager.AddOrchestrator(orchestrator)
	pluginOrchestrator, err := manager.Get("orchestrator")
	assert.Nil(t, err)
	data, err := pluginOrchestrator.Execute(context.Background(), map[string]any{"testPlugin": true})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test", "testPlugin": true}, data)
	data, err = pluginOrchestrator.Execute(context.Background(), map[string]any{"secret": "secret"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"secret": "c2VjcmV0"}, data)
}

func TestOrchestratorPluginFailureCases(t *testing.T) {
	manager := plugins.NewBasePluginManager()
	orchestrator := plugins.NewTransformPlugin("orchestrator", func(data any) (any, error) {
		dataMap, ok := data.(map[string]any)
		if !ok {
			return nil, errors.New("data is not a map")
		}
		if dataMap["test"] != nil {
			return "test", nil
		} else {
			return nil, errors.New("invalid input")
		}
	})

	manager.AddOrchestrator(orchestrator)
	pluginOrchestrator, err := manager.Get(orchestrator.GetName())
	assert.Nil(t, err)
	data, err := pluginOrchestrator.Execute(context.Background(), map[string]any{"test": "test"})
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "plugin not found")

	data, err = pluginOrchestrator.Execute(context.Background(), map[string]any{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "invalid input")

	manager.AddOrchestrator(testPlugin)
	pluginOrchestrator, err = manager.Get(testPlugin.GetName())
	assert.Nil(t, err)
	data, err = pluginOrchestrator.Execute(context.Background(), map[string]any{"test": "test"})
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "plugin is not an orchestrator")
}
