package plugins_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/stretchr/testify/assert"
)

const (
	someError  = "some error"
	helloWorld = "Hello World!!"
)

func createWorkflowFromFile(workflowFile string, pluginManager plugins.PluginManager) (*plugins.WorkflowConfig, plugins.WorkflowPlugin, error) {
	workflowConfig, err := plugins.LoadWorkflowFile(workflowFile)
	if err != nil {
		return nil, nil, err
	}
	workflowPlugin, err := plugins.NewBaseWorkflowPlugin(pluginManager, *workflowConfig)
	return workflowConfig, workflowPlugin, err
}

func createSampleWorkflow() (*plugins.WorkflowConfig, plugins.WorkflowPlugin, error) {
	pluginManager := plugins.NewBasePluginManager()
	pluginManager.Add(testPlugin)
	pluginManager.Add(bloblPlugin)
	pluginManager.Add(badPlugin)

	return createWorkflowFromFile("../test_data/workflows/sample.yaml", pluginManager)
}

func TestWorkflowEngine(t *testing.T) {
	sampleWorkflowConfig, workflowPlugin, err := createSampleWorkflow()
	assert.Nil(t, err)
	assert.NotNil(t, workflowPlugin)
	assert.Equal(t, sampleWorkflowConfig.Name, workflowPlugin.GetName())
	bloblStep, err := workflowPlugin.GetStep("blobl1")
	assert.Nil(t, err)
	assert.NotNil(t, bloblStep)
	assert.Equal(t, "blobl1", bloblStep.GetName())
	assert.Equal(t, plugins.BloblangStep, bloblStep.GetType())
	steps := workflowPlugin.GetSteps()
	assert.Equal(t, len(sampleWorkflowConfig.Steps), len(steps))

	result, err := workflowPlugin.ExecuteStep(context.Background(), "blobl1", emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, "Hello", result.Data)

	result, err = workflowPlugin.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, result.Data)

	result, err = workflowPlugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, helloWorld, result.Data)
}

func TestWorkflowEngineReplay(t *testing.T) {
	workflowConfig, workflowPlugin, err := createSampleWorkflow()
	assert.Nil(t, err)
	input := emptyMessage()
	input.Version = workflowConfig.Version
	input.Status.LastCompletedStepIndex = 4
	input.SetMetadata(plugins.WorkflowOutputsKey, map[string]any{"blobl2": "Hello World"})
	result, err := workflowPlugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	assert.Equal(t, helloWorld, result.Data)

	input.Status.LastCompletedStepIndex = 5
	result, err = workflowPlugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	// No step will be executed as the last completed step index is
	// greater than or equals last workflow step index
	assert.Equal(t, input.Data, result.Data)

	input = testMessage()
	// Set version to previous version to test replay
	// Now the last completed step index will be reset to 0
	// Workflow will be executed from the beginning
	input.Version = workflowConfig.Version - 1
	input.Status.LastCompletedStepIndex = 3
	result, err = workflowPlugin.Execute(context.Background(), input)
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"test": "test"}, result.Data)
	assert.Equal(t, workflowConfig.Version, result.Version)

	pluginsManager := plugins.NewBasePluginManager()
	pluginsManager.Add(newFailingPlugin("unreliable", someError, 1))

	_, unreliableWorkflowPlugin, err := createWorkflowFromFile("../test_data/workflows/unreliable.yaml", pluginsManager)
	assert.Nil(t, err)
	workflowStatus := plugins.ExecutionStatus{}
	assert.True(t, workflowStatus.IsUnprocessed())
	assert.Equal(t, plugins.ExecutionStatusUnprocessed, workflowStatus.GetStatus())

	result, err = unreliableWorkflowPlugin.Execute(context.Background(), emptyMessage())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, someError)
	assert.Equal(t, "bar", result.Data)

	assert.Equal(t, plugins.ExecutionStatusFailed, result.Status.GetStatus())
	assert.Equal(t, someError, result.Status.Message)
	assert.Equal(t, 0, result.Status.LastCompletedStepIndex)

	result, err = unreliableWorkflowPlugin.Execute(context.Background(), result)
	assert.Nil(t, err)

	assert.Equal(t, len(unreliableWorkflowPlugin.GetSteps())-1, result.Status.LastCompletedStepIndex)
	assert.Equal(t, plugins.ExecutionStatusCompleted, result.Status.Status)
	assert.True(t, result.Status.IsCompleted())
	assert.Equal(t, helloWorld, result.Data)
}

func TestStepPlugin(t *testing.T) {
	pluginsManager := plugins.NewBasePluginManager()
	pluginsManager.Add(plugins.NewBasePlugin("test", newFailingExecutor(someError, 1)))

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
						Name:       "test",
						CheckBlobl: "some check",
						Bloblang:   "this",
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
						Name:      "test",
						CheckExpr: `some error`,
						Bloblang:  `this`,
					},
				},
			},
			expectedError: "failed to parse expression",
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name: "test",
						Expr: `some error`,
					},
				},
			},
			expectedError: "failed to parse expression",
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
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:       "test",
						CheckBlobl: "this",
						CheckExpr:  "this",
						Bloblang:   "this",
					},
				},
			},
			expectedError: "only one of check_blobl and check_expr is allowed",
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
		executionError string
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
			executionError: someError,
		},
		{
			workflowConfig: plugins.WorkflowConfig{
				Name: "test",
				Steps: []plugins.StepConfig{
					{
						Name:       "test",
						CheckBlobl: `throw("some error")`,
						Bloblang:   `this`,
					},
				},
			},
			executionError: someError,
		},
	}
	for _, testCase := range testCases {
		pluginManager := plugins.NewBasePluginManager()
		workflowPlugin, err := plugins.NewBaseWorkflowPlugin(pluginManager, testCase.workflowConfig)
		assert.Nil(t, err)
		assert.NotNil(t, workflowPlugin)
		result, err := workflowPlugin.Execute(context.Background(), emptyMessage())
		assert.NotNil(t, result)
		assert.NotNil(t, err)
		assert.True(t, result.Status.IsFailed())
		assert.ErrorContains(t, err, testCase.executionError)
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

	workflow, err = plugins.LoadWorkflowFile("workflow.go")
	assert.Nil(t, workflow)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to unmarshal workflow file")
}
