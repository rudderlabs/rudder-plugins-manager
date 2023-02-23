package destinations

import (
	"errors"

	"github.com/rudderlabs/rudder-transformations/plugins/types"
)

type DefaultPlugin struct{}

func (p *DefaultPlugin) Name() string {
	return "default"
}

func (p *DefaultPlugin) GetTransformer(data interface{}) (types.Transformer, error) {
	return p, nil
}

func (p *DefaultPlugin) Transform(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["default"] = true
	return dataMap, nil
}
