package plugins

import (
	"context"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/rudderlabs/rudder-plugins-manager/types"
)

/**
 * This is a base plugin that can be used to create a plugin from a function.
 */
type BasePlugin struct {
	Name     string
	Executor types.Executor
}

func NewBasePlugin(name string, executor types.Executor) *BasePlugin {
	return &BasePlugin{Name: name, Executor: executor}
}

func (p *BasePlugin) Execute(ctx context.Context, data any) (any, error) {
	return p.Executor.Execute(ctx, data)
}

func (p *BasePlugin) GetName() string {
	return p.Name
}

func NewTransformPlugin(name string, fn func(any) (any, error)) types.Plugin {
	return &BasePlugin{
		Name:     name,
		Executor: types.TransformFunc(fn),
	}
}

/**
 * This is transforms the data using bloblang template.
 */
type BloblangPlugin struct {
	Name     string
	executor *bloblang.Executor
}

func NewBloblangPlugin(name, template string) (*BloblangPlugin, error) {
	executor, err := bloblang.Parse(template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bloblang template: %w", err)
	}
	return &BloblangPlugin{Name: name, executor: executor}, nil
}

func (p *BloblangPlugin) Execute(_ context.Context, input any) (any, error) {
	return p.executor.Query(input)
}

func (p *BloblangPlugin) GetName() string {
	return p.Name
}

type OrchestratorPlugin struct {
	manager types.PluginManager
	plugin  types.Plugin
}

func NewOrchestratorPlugin(manager types.PluginManager, plugin types.Plugin) *OrchestratorPlugin {
	return &OrchestratorPlugin{manager: manager, plugin: plugin}
}

func (p *OrchestratorPlugin) Execute(ctx context.Context, data any) (any, error) {
	result, err := p.plugin.Execute(ctx, data)
	if err != nil {
		return nil, err
	}
	pluginName, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("plugin is not an orchestrator: result must be a string")
	}
	nextPlugin, err := p.manager.GetPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %s, error: %w", pluginName, err)
	}
	return nextPlugin.Execute(ctx, data)
}

func (p *OrchestratorPlugin) GetName() string {
	return p.plugin.GetName()
}
