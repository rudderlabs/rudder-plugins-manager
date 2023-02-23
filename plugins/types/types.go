package types

/**
 * This file contains the interface that all plugin providers must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type PluginProvider interface {
	Name() string
	AddPlugin(Plugin)
	GetPlugin(string) (Plugin, bool)
}

/**
 * This file contains the interface that all plugins must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type Plugin interface {
	Name() string
	GetTransformer(interface{}) (Transformer, error)
}

/**
 * This interface is used by the plugin manager to transform data.
 * The plugin manager will call GetTransformer() on the plugin to get a transformer.
 * The plugin manager will then call Transform() on the transformer to transform the data.
 */
type Transformer interface {
	Transform(interface{}) (interface{}, error)
}

/**
 * This is a helper function that allows you to create a Transformer from a function.
 */
type TransformerFunc func(interface{}) (interface{}, error)

/**
 * This is a helper function that allows you to create a Transformer from a function.
 */
func (f TransformerFunc) Transform(data interface{}) (interface{}, error) {
	return f(data)
}
