package plugins

import (
	"context"
	"errors"
)

/**
 * This file contains the interface that all plugin providers must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type BasePluginManager struct {
	plugins map[string]Plugin
}

func NewBasePluginManager() *BasePluginManager {
	return &BasePluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (m *BasePluginManager) Add(plugin Plugin) {
	m.plugins[plugin.GetName()] = plugin
}

func (m *BasePluginManager) AddOrchestrator(plugin Plugin) {
	m.Add(&OrchestratorPlugin{
		manager: m,
		plugin:  plugin,
	})
}

func (p *BasePluginManager) Get(name string) (Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, errors.New("plugin not found")
	}
	return plugin, nil
}

func (p *BasePluginManager) Execute(ctx context.Context, name string, data *Message) (*Message, error) {
	plugin, err := p.Get(name)
	if err != nil {
		return nil, err
	}
	return plugin.Execute(ctx, data)
}
