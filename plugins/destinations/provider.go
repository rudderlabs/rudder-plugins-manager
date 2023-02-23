package destinations

import (
	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

const PROVIDER_NAME = "destinations"

/**
 * This struct is used to manage plugins.
 */
type PluginProvider struct {
	plugins map[string]types.Plugin
}

var ProviderInstance types.PluginProvider

/**
 * This function creates a new plugin Provider.
 */
func NewPluginProvider() *PluginProvider {
	return &PluginProvider{
		plugins: map[string]types.Plugin{
			"default": &DefaultPlugin{},
		},
	}
}

/**
 * This function returns the name of the plugin Provider.
 */
func (p *PluginProvider) Name() string {
	return PROVIDER_NAME
}

/**
 * This function adds a plugin to the plugin Provider.
 */
func (p *PluginProvider) AddPlugin(plugin types.Plugin) {
	p.plugins[plugin.Name()] = plugin
}

/**
 * This function gets a plugin from the plugin Provider.
 */
func (p *PluginProvider) GetPlugin(name string) (types.Plugin, bool) {
	plugin, ok := p.plugins[name]
	return plugin, ok
}

/**
 * This function initializes the plugin Provider.
 */
func init() {
	ProviderInstance = NewPluginProvider()
	ProviderInstance.AddPlugin(&DefaultPlugin{})
}
