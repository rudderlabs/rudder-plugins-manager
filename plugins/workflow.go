package plugins

import (
	"context"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/mitchellh/mapstructure"
)

type StepInput struct {
	Data    any            `mapstructure:"data"`
	Input   any            `mapstructure:"input"`
	Context map[string]any `mapstructure:"context"`
}

func (s *StepInput) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data":    s.Data,
		"input":   s.Input,
		"context": s.Context,
	}
}

type StepOutput struct {
	Output  any            `mapstructure:"output"`
	Context map[string]any `mapstructure:"context"`
}

type StepPlugin interface {
	GetName() string
	ShouldExecute(data StepInput) (bool, error)
	Execute(ctx context.Context, data StepInput) (*StepOutput, error)
}

type BaseStepPlugin struct {
	Name          string
	ConditionFunc func(data StepInput) (bool, error)
	ExecuteFunc   func(ctx context.Context, data StepInput) (*StepOutput, error)
}

func (p *BaseStepPlugin) GetName() string {
	return p.Name
}

func (p *BaseStepPlugin) ShouldExecute(data StepInput) (bool, error) {
	if p.ConditionFunc != nil {
		return p.ConditionFunc(data)
	}
	return true, nil
}

func (p *BaseStepPlugin) Execute(ctx context.Context, data StepInput) (*StepOutput, error) {
	if p.ExecuteFunc != nil {
		return p.ExecuteFunc(ctx, data)
	}
	return nil, nil
}

type WorkflowPlugin struct {
	Name  string
	Steps []StepPlugin
}

func NewWorkflowPlugin(name string, steps ...StepPlugin) *WorkflowPlugin {
	return &WorkflowPlugin{
		Name:  name,
		Steps: steps,
	}
}

func (w *WorkflowPlugin) GetName() string {
	return w.Name
}

func (w *WorkflowPlugin) Execute(ctx context.Context, input any) (any, error) {
	stepInput := StepInput{
		Data:    input,
		Input:   input,
		Context: make(map[string]any),
	}

	for _, step := range w.Steps {
		if shouldExecute, err := step.ShouldExecute(stepInput); err != nil {
			return nil, err
		} else if shouldExecute {
			if stepOutput, err := step.Execute(ctx, stepInput); err != nil {
				return nil, err
			} else {
				stepInput.Data = stepOutput.Output
				stepInput.Context = stepOutput.Context
			}
		}
	}
	return stepInput.Data, nil
}

type BloblangWorkflow struct {
	Name  string                 `yaml:"name"`
	Steps []BloblangWorkflowStep `yaml:"steps"`
}

func NewBloblangWorkflowPlugin(workflow *BloblangWorkflow) (*WorkflowPlugin, error) {
	var steps []StepPlugin
	for _, step := range workflow.Steps {
		if plugin, err := NewBloblangStepPlugin(step); err != nil {
			return nil, err
		} else {
			steps = append(steps, plugin)
		}
	}
	return &WorkflowPlugin{
		Name:  workflow.Name,
		Steps: steps,
	}, nil
}

type BloblangWorkflowStep struct {
	Name      string  `yaml:"name"`
	Condition *string `yaml:"condition"`
	Template  string  `yaml:"template"`
}

type BloblangStepPlugin struct {
	Name      string
	condition *bloblang.Executor
	template  *bloblang.Executor
}

func NewBloblangStepPlugin(step BloblangWorkflowStep) (*BloblangStepPlugin, error) {
	stepPlugin := &BloblangStepPlugin{
		Name: step.Name,
	}
	if step.Condition != nil {
		if condition, err := bloblang.Parse(*step.Condition); err != nil {
			return nil, fmt.Errorf("failed to parse condition for step %s: %w", step.Name, err)
		} else {
			stepPlugin.condition = condition
		}
	}
	if template, err := bloblang.Parse(step.Template); err != nil {
		return nil, fmt.Errorf("failed to parse template for step %s: %w", step.Name, err)
	} else {
		stepPlugin.template = template
	}

	return stepPlugin, nil
}

func (b *BloblangStepPlugin) GetName() string {
	return b.Name
}

func (b *BloblangStepPlugin) ShouldExecute(input StepInput) (bool, error) {
	if b.condition == nil {
		return true, nil
	}

	if result, err := b.condition.Query(input.ToMap()); err != nil {
		return false, err
	} else {
		return result.(bool), nil
	}
}

func (b *BloblangStepPlugin) Execute(_ context.Context, input StepInput) (*StepOutput, error) {
	var output StepOutput
	if result, err := b.template.Query(input.ToMap()); err != nil {
		return nil, err
	} else {
		if err := mapstructure.Decode(result, &output); err != nil {
			return nil, err
		}
		return &output, nil
	}
}
