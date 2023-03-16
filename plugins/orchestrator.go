package plugins

import (
	"context"
	"fmt"
)

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
