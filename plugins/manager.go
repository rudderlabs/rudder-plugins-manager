package plugins

import (
	"errors"

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
