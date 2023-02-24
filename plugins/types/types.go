package types

import "errors"

/**
 * This file contains the interface that all plugin providers must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type PluginProvider struct {
	Name    string
	plugins map[string]Plugin
}

func NewPluginProvider(name string) *PluginProvider {
	return &PluginProvider{
		Name:    name,
		plugins: make(map[string]Plugin),
	}
}

func (p *PluginProvider) GetName() string {
	return p.Name
}

func (p *PluginProvider) AddPlugin(plugin Plugin) {
	p.plugins[plugin.GetName()] = plugin
}

func (p *PluginProvider) GetPlugin(name string) (Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, errors.New("plugin not found")
	}
	return plugin, nil
}

/**
 * This file contains the interface that all plugins must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type Plugin interface {
	GetName() string
	GetTransformer(any) (Transformer, error)
}

/**
 * This interface is used by the plugin manager to transform data.
 * The plugin manager will call GetTransformer() on the plugin to get a transformer.
 * The plugin manager will then call Transform() on the transformer to transform the data.
 */
type Transformer interface {
	Transform(any) (any, error)
}

/**
 * This is a helper function that allows you to create a Transformer from a function.
 */
type TransformerFunc func(any) (any, error)

/**
 * This is a helper function that allows you to create a Transformer from a function.
 */
func (f TransformerFunc) Transform(data any) (any, error) {
	return f(data)
}
