package plugins

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/rs/zerolog/log"
)

/**
 * This is a base plugin that can be used to create a plugin from a function.
 */
type BasePlugin struct {
	Name     string
	Executor Executor
}

func NewBasePlugin(name string, executor Executor) *BasePlugin {
	return &BasePlugin{Name: name, Executor: executor}
}

func (p *BasePlugin) Execute(ctx context.Context, data *Message) (*Message, error) {
	return p.Executor.Execute(ctx, data)
}

func (p *BasePlugin) GetName() string {
	return p.Name
}

func NewTransformPlugin(name string, fn func(*Message) (*Message, error)) Plugin {
	return &BasePlugin{
		Name:     name,
		Executor: TransformFunc(fn),
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

func (p *BloblangPlugin) Execute(_ context.Context, input *Message) (*Message, error) {
	inputMap := input.ToMap()
	log.Debug().Any("input", inputMap).Msg("Executing bloblang plugin")
	data, err := p.executor.Query(inputMap)
	if err != nil {
		return nil, err
	}
	var result Message
	var dataBytes []byte
	if dataBytes, err = json.Marshal(data); err == nil {
		err = json.Unmarshal(dataBytes, &result)
	}
	if err != nil || (result.Metadata == nil && result.Data == nil) {
		return NewMessage(data), nil
	}
	return &result, nil
}

func (p *BloblangPlugin) GetName() string {
	return p.Name
}

type OrchestratorPlugin struct {
	manager PluginManager
	plugin  Plugin
}

func NewOrchestratorPlugin(manager PluginManager, plugin Plugin) *OrchestratorPlugin {
	return &OrchestratorPlugin{manager: manager, plugin: plugin}
}

func NextPluginMessage(nextPlugin string) *Message {
	return &Message{
		Metadata: map[string]any{
			"next_plugin": nextPlugin,
		},
	}
}

func (p *OrchestratorPlugin) Execute(ctx context.Context, data *Message) (*Message, error) {
	result, err := p.plugin.Execute(ctx, data)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	pluginName, ok := result.Metadata["next_plugin"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get next plugin name")
	}
	nextPlugin, err := p.manager.Get(pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %s, error: %w", pluginName, err)
	}
	return nextPlugin.Execute(ctx, data)
}

func (p *OrchestratorPlugin) GetName() string {
	return p.plugin.GetName()
}
