package destinations

import (
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
	return data, nil
}
