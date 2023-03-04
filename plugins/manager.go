package plugins

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/rudderlabs/rudder-plugins-manager/types"
)

/**
 * This file contains the interface that all plugin providers must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type PluginManager struct {
	plugins map[string]types.Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]types.Plugin),
	}
}

func (p *PluginManager) AddPlugin(plugin types.Plugin) {
	p.plugins[plugin.GetName()] = plugin
}

func (p *PluginManager) GetPlugin(name string) (types.Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, errors.New("plugin not found")
	}
	return plugin, nil
}

func (p *PluginManager) Execute(ctx context.Context, name string, data any) (any, error) {
	plugin, err := p.GetPlugin(name)
	if err != nil {
		return nil, err
	}
	result, err := plugin.Execute(ctx, data)
	if err != nil {
		return nil, err
	}
	var nextPlugin types.NextPlugin
	if err := mapstructure.Decode(result, &nextPlugin); err != nil {
		return result, nil
	}
	if nextPlugin.NextPluginName != nil {
		return p.Execute(ctx, *nextPlugin.NextPluginName, nextPlugin.Data)
	}
	return result, nil
}
