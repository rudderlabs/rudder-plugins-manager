package plugins_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowEngine(t *testing.T) {
	pluginManager := plugins.NewBasePluginManager()
	pluginManager.Add(testPlugin)
	pluginManager.Add(bloblPlugin)
	pluginManager.Add(badPlugin)

	workflowConfig1, err := plugins.LoadWorkflowFile("../test_data/workflows/workflow1.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, workflowConfig1)
	workflowPlugin, err := plugins.NewBaseWorkflowPlugin(pluginManager, *workflowConfig1)
	assert.Nil(t, err)
	assert.NotNil(t, workflowPlugin)
	assert.Equal(t, workflowConfig1.Name, workflowPlugin.GetName())
	bloblStep, err := workflowPlugin.GetStep("blobl")
	assert.Nil(t, err)
	assert.NotNil(t, bloblStep)
	assert.Equal(t, "blobl", bloblStep.GetName())
	assert.Equal(t, plugins.BloblangStep, bloblStep.GetType())
	steps := workflowPlugin.GetSteps()
	assert.Equal(t, len(workflowConfig1.Steps), len(steps))

	result, err := workflowPlugin.ExecuteStep(context.Background(), "blobl", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"blobl": true}, result.Data)

	result, err = workflowPlugin.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, result.Data)

	result, err = workflowPlugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"blobl": true}, result.Data)
}

func TestStepPlugin(t *testing.T) {
	pluginsManager := plugins.NewBasePluginManager()
	pluginsManager.Add(plugins.NewBasePlugin("test", newFailingExecutor("some error", 1)))

	stepPlugin, err := plugins.NewBaseStepPlugin(pluginsManager, plugins.StepConfig{
		Name:   "test",
		Plugin: "test",
		Retry: &plugins.BaseRetryPolicy{
			RetryCount: 1,
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, stepPlugin)
	assert.Equal(t, "test", stepPlugin.GetName())
	assert.Equal(t, plugins.PluginStep, stepPlugin.GetType())
	policy, ok := stepPlugin.GetRetryPolicy()
	assert.True(t, ok)
	assert.Equal(t, uint64(1), policy.GetRetryCount())
	result, err := stepPlugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, emptyMessage(), result)

	stepPlugin, err = plugins.NewBaseStepPlugin(pluginsManager, plugins.StepConfig{
		Name:   "test",
		Plugin: "test",
	})
	assert.Nil(t, err)
	_, ok = stepPlugin.GetRetryPolicy()
	assert.False(t, ok)
}

func TestWorkflowInvalidConfig(t *testing.T) {
	type errorTestCase struct {
		workflowConfig plugins.WorkflowConfig
		expectedError  string
		pluginManager  plugins.PluginManager
	}

	testCases := []errorTestCase{
		{
			expectedError: "workflow name is required",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
			},
			expectedError: "workflow must have at least one step",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{},
				},
			},
			expectedError: "step name is required",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name: "test",
					},
				},
			},
			expectedError: "unknown step type",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:   "test",
						Plugin: "test",
					},
				},
			},
			expectedError: "plugin manager is required when plugin is set",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Bloblang: "this",
						Return:   true,
					},
				},
			},
			expectedError: "return is only allowed when check is set",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Check:    "some check",
						Bloblang: "this",
					},
				},
			},
			expectedError: "failed to parse bloblang template",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Bloblang: "some bloblang",
					},
				},
			},
			expectedError: "failed to parse bloblang template",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:   "test",
						Plugin: "test",
					},
				},
			},
			expectedError: "plugin not found",
			pluginManager: plugins.NewBasePluginManager(),
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Bloblang: "this",
					},
					{
						Name:     "test",
						Bloblang: "this",
					},
				},
			},
			expectedError: "workflow steps must have unique names",
		},
	}

	for _, testCase := range testCases {
		workflowPlugin, err := plugins.NewBaseWorkflowPlugin(testCase.pluginManager, testCase.workflowConfig)
		assert.Nil(t, workflowPlugin)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, testCase.expectedError)
	}
}

func TestWorkflowExecutionFailures(t *testing.T) {
	type errorTestCase struct {
		workflowConfig plugins.WorkflowConfig
		expectedError  string
	}

	testCases := []errorTestCase{
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Bloblang: `throw("some error")`,
					},
				},
			},
			expectedError: "some error",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:     "test",
						Check:    `throw("some error")`,
						Bloblang: `this`,
					},
				},
			},
			expectedError: "some error",
		},
	}
	for _, testCase := range testCases {
		pluginManager := plugins.NewBasePluginManager()
		workflowPlugin, err := plugins.NewBaseWorkflowPlugin(pluginManager, testCase.workflowConfig)
		assert.Nil(t, err)
		assert.NotNil(t, workflowPlugin)
		result, err := workflowPlugin.Execute(context.Background(), emptyMessage())
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, testCase.expectedError)
	}
}

func TestWorkflowGetStepFailureCase(t *testing.T) {
	workflowPlugin, err := plugins.NewBaseWorkflowPlugin(nil, plugins.WorkflowConfig{
		Name: "test",
		Steps: []plugins.StepConfig{
			{
				Name:     "test",
				Bloblang: "this",
			},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, workflowPlugin)
	_, err = workflowPlugin.GetStep("not-existent")
	assert.ErrorContains(t, err, "step not-existent not found")
}

func TestWorkflowExecuteStepFailureCases(t *testing.T) {
	workflowPlugin, err := plugins.NewBaseWorkflowPlugin(nil, plugins.WorkflowConfig{
		Name: "test",
		Steps: []plugins.StepConfig{
			{
				Name:     "test",
				Bloblang: "this",
			},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, workflowPlugin)
	_, err = workflowPlugin.ExecuteStep(context.Background(), "not-existent", emptyMessage())
	assert.ErrorContains(t, err, "step not-existent not found")
}

func TestLoadWorkflowFileFailureCases(t *testing.T) {
	workflow, err := plugins.LoadWorkflowFile("non-existent")
	assert.Nil(t, workflow)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to read workflow file")

	workflow, err = plugins.LoadWorkflowFile("workflows.go")
	assert.Nil(t, workflow)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to unmarshal workflow file")
}
