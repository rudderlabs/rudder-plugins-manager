package destinations

import (
	"errors"

	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

const PROVIDER_NAME = "destinations"

var ProviderInstance *types.PluginProvider

var DefaultPlugin = types.NewSimplePlugin(
	"default",
	func(data interface{}) (interface{}, error) {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return nil, errors.New("data is not a map")
		}
		dataMap["default"] = true
		return dataMap, nil
	},
)

/**
 * This function initializes the plugin Provider.
 */
func init() {
	ProviderInstance = types.NewPluginProvider(PROVIDER_NAME)
}
