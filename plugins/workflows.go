package plugins

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type BaseStepPlugin struct {
	Name     string
	Type     StepType
	Check    Executor
	Main     Executor
	Continue bool
	Return   bool
}

func validateStepConfig(config *StepConfig, pluginManager PluginManager) error {
	if config.Name == "" {
		return fmt.Errorf("step name is required")
	}

	if config.GetType() == UnknownStep {
		return fmt.Errorf("unknown step type")
	}

	if config.Plugin != "" && pluginManager == nil {
		return fmt.Errorf("plugin manager is required when plugin is set")
	}

	if config.Check == "" && config.Return {
		return fmt.Errorf("return is only allowed when check is set")
	}
	return nil
}

func NewBaseStepPlugin(pluginManager PluginManager, config StepConfig) (StepPlugin, error) {
	err := validateStepConfig(&config, pluginManager)
	if err != nil {
		return nil, err
	}
	stepPlugin := BaseStepPlugin{
		Name:     config.Name,
		Continue: config.Continue,
		Return:   config.Return,
	}

	stepPlugin.Type = config.GetType()

	switch stepPlugin.Type {
	case BloblangStep:
		executor, err := NewBloblangPlugin(config.Name, config.Bloblang)
		if err != nil {
			return nil, err
		}
		stepPlugin.Main = executor
	default:
		plugin, err := pluginManager.Get(config.Plugin)
		if err != nil {
			return nil, err
		}
		stepPlugin.Main = plugin
	}

	if config.Check != "" {
		check, err := NewBloblangPlugin(fmt.Sprintf("%s.check", config.Name), config.Check)
		if err != nil {
			return nil, err
		}
		stepPlugin.Check = check
	}
	return &stepPlugin, nil
}

func (p *BaseStepPlugin) GetType() StepType {
	return p.Type
}

func (p *BaseStepPlugin) ShouldExecute(data *Message) (bool, error) {
	if p.Check == nil {
		return true, nil
	}
	output, err := p.Check.Execute(context.Background(), data)
	if err != nil {
		return false, err
	}
	return output.GetBool()
}

func (p *BaseStepPlugin) ShouldReturn() bool {
	return p.Return
}

func (p *BaseStepPlugin) ShouldContinue() bool {
	return p.Continue
}

func (p *BaseStepPlugin) GetName() string {
	return p.Name
}

func (p *BaseStepPlugin) Execute(ctx context.Context, data *Message) (*Message, error) {
	return p.Main.Execute(ctx, data)
}

type BaseWorkflowPlugin struct {
	Name  string
	Steps []StepPlugin
}

func validateWorkflowConfig(config *WorkflowConfig) error {
	if config.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(config.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	uniqSteps := lo.UniqBy(config.Steps, func(step StepConfig) string {
		return step.Name
	})

	if len(uniqSteps) != len(config.Steps) {
		return fmt.Errorf("workflow steps must have unique names")
	}
	return nil
}

func NewBaseWorkflowPlugin(pluginManager PluginManager, config WorkflowConfig) (WorkflowPlugin, error) {
	err := validateWorkflowConfig(&config)
	if err != nil {
		return nil, err
	}
	steps := make([]StepPlugin, len(config.Steps))
	for idx, stepConfig := range config.Steps {
		steps[idx], err = NewBaseStepPlugin(pluginManager, stepConfig)
		if err != nil {
			return nil, err
		}
	}

	return &BaseWorkflowPlugin{Name: config.Name, Steps: steps}, nil
}

func (p *BaseWorkflowPlugin) GetName() string {
	return p.Name
}

func (p *BaseWorkflowPlugin) Execute(ctx context.Context, input *Message) (*Message, error) {
	newInput := input.Clone()
	for _, step := range p.Steps {
		shouldExecute, err := step.ShouldExecute(newInput)
		if err != nil {
			if step.ShouldContinue() {
				log.Warn().Err(err).Msg("failed to check step")
				continue
			}
			return nil, err
		}
		if shouldExecute {
			output, err := step.Execute(ctx, newInput)
			if err != nil {
				if step.ShouldContinue() {
					log.Warn().Err(err).Msg("failed to execute step")
					continue
				}
				return nil, err
			}

			newInput.Data = output.Data
			newInput.Metadata = lo.Assign(newInput.Metadata, output.Metadata)
			if step.ShouldReturn() {
				return newInput, nil
			}
		}
	}
	return newInput, nil
}

func (p *BaseWorkflowPlugin) GetSteps() []StepPlugin {
	return p.Steps
}

func (p *BaseWorkflowPlugin) GetStep(name string) (StepPlugin, error) {
	for _, step := range p.Steps {
		if step.GetName() == name {
			return step, nil
		}
	}
	return nil, fmt.Errorf("step %s not found", name)
}

func (p *BaseWorkflowPlugin) ExecuteStep(ctx context.Context, stepName string, data *Message) (*Message, error) {
	step, err := p.GetStep(stepName)
	if err != nil {
		return nil, err
	}
	return step.Execute(ctx, data)
}
