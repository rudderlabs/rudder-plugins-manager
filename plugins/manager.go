package plugins

import (
	"errors"

	"github.com/rudderlabs/rudder-transformations/plugins/destinations"
	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

/**
 * This struct is used to manage plugins.
 */
type Manager struct {
	providers map[string]*types.PluginProvider
}

var ManagerInstance *Manager

/**
 * This function creates a new plugin manager.
 */
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]*types.PluginProvider),
	}
}

/**
 * This function adds a plugin to the plugin manager.
 */
func (m *Manager) AddPluginProvider(provider *types.PluginProvider) {
	m.providers[provider.GetName()] = provider
}

/**
 * This function adds a plugin to the plugin manager.
 */
func (m *Manager) AddPlugin(provider string, plugin types.Plugin) error {
	pluginProvider, ok := m.providers[provider]
	if !ok {
		return errors.New("provider not found")
	}
	pluginProvider.AddPlugin(plugin)
	return nil
}

func (m *Manager) AddPlugins(provider string, plugins ...types.Plugin) error {
	for _, plugin := range plugins {
		if err := m.AddPlugin(provider, plugin); err != nil {
			return err
		}
	}
	return nil
}

/**
 * This function gets a plugin from the plugin manager.
 */
func (m *Manager) GetPlugin(provider, plugin string) (types.Plugin, error) {
	pluginProvider, ok := m.providers[provider]
	if !ok {
		return nil, errors.New("provider not found")
	}
	return pluginProvider.GetPlugin(plugin)
}

/**
 * This function initializes the plugin manager.
 */
func init() {
	ManagerInstance = NewManager()
	ManagerInstance.AddPluginProvider(destinations.ProviderInstance)
}
