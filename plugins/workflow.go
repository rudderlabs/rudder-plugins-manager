package plugins

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const (
	WorkflowOutputsKey = "outputs"
)

type BaseStepPlugin struct {
	Name        string
	Type        StepType
	Check       Executor
	Main        Executor
	Continue    bool
	Return      bool
	RetryPolicy RetryPolicy
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

func getMainStepExecutor(config *StepConfig, pluginManager PluginManager) (Executor, error) {
	switch config.GetType() {
	case BloblangStep:
		return NewBloblangPlugin(config.Name, config.Bloblang)
	default:
		plugin, err := pluginManager.Get(config.Plugin)
		if err != nil {
			return nil, err
		}
		return plugin, nil
	}
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

	stepPlugin.Main, err = getMainStepExecutor(&config, pluginManager)
	if err != nil {
		return nil, err
	}

	if config.Retry != nil {
		stepPlugin.RetryPolicy = config.Retry
		stepPlugin.Main = NewBaseRetryableExecutor(stepPlugin.Main, stepPlugin.RetryPolicy)
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

func (p *BaseStepPlugin) GetRetryPolicy() (RetryPolicy, bool) {
	if p.RetryPolicy == nil {
		return nil, false
	}
	return p.RetryPolicy, true
}

type BaseWorkflowPlugin struct {
	Name    string
	Version int
	Steps   []StepPlugin
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

	return &BaseWorkflowPlugin{Name: config.Name, Version: config.GetVersion(), Steps: steps}, nil
}

func (p *BaseWorkflowPlugin) GetName() string {
	return p.Name
}

func (p *BaseWorkflowPlugin) GetVersion() int {
	return p.Version
}

func executeWorkflowStep(ctx context.Context, step StepPlugin, data *Message) (*Message, error) {
	shouldExecute, err := step.ShouldExecute(data)
	if err != nil {
		if step.ShouldContinue() {
			log.Warn().Err(err).Msg("failed to check step")
			return nil, nil
		}
		return nil, err
	}
	if !shouldExecute {
		return nil, nil
	}

	output, err := step.Execute(ctx, data)
	if err != nil {
		if step.ShouldContinue() {
			log.Warn().Err(err).Msg("failed to execute step")
			return nil, nil
		}
		return nil, err
	}
	return output, nil
}

func initWorkflowMessage(workflow WorkflowPlugin, input *Message) *Message {
	newInput := input.Clone()
	if newInput.Version != workflow.GetVersion() {
		newInput.Version = workflow.GetVersion()
		newInput.Status.LastCompletedStepIndex = -1
		newInput.Metadata[WorkflowOutputsKey] = map[string]any{}
	}
	return newInput
}

func (p *BaseWorkflowPlugin) Execute(ctx context.Context, input *Message) (*Message, error) {
	newInput := initWorkflowMessage(p, input)
	startIdx := newInput.Status.LastCompletedStepIndex + 1
	log.Debug().Str("workflow", p.Name).Int("startIdx", startIdx).Msg("Execution is started")
	for i := startIdx; i < len(p.Steps); i++ {
		step := p.Steps[i]
		output, err := executeWorkflowStep(ctx, step, newInput)
		if err != nil {
			newInput.Status.SetError(err)
			return newInput, err
		}
		newInput.Status.LastCompletedStepIndex = i
		if output == nil {
			continue
		}
		newInput.Data = output.Data
		newInput.Metadata = lo.Assign(newInput.Metadata, output.Metadata)
		newInput.Metadata[WorkflowOutputsKey].(map[string]any)[step.GetName()] = output.Data
		if step.ShouldReturn() {
			return newInput, nil
		}
	}
	newInput.Status.Status = ExecutionStatusCompleted
	log.Debug().Str("workflow", p.Name).Msg("Execution is successful")
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

func LoadWorkflowFile(workflowFile string) (*WorkflowConfig, error) {
	var workflowConfig WorkflowConfig
	workflowBytes, err := os.ReadFile(workflowFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}
	err = yaml.Unmarshal(workflowBytes, &workflowConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow file: %w", err)
	}
	return &workflowConfig, nil
}
