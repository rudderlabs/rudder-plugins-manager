package types_test

import (
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-transformations/plugins/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSimplePlugin(t *testing.T) {
	plugin := types.NewSimplePlugin("test", func(data any) (any, error) {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return nil, errors.New("data is not a map")
		}
		dataMap["test"] = "test"
		return dataMap, nil
	})
	data, err := plugin.Transform(map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"test": "test"}, data)
}

func TestNewSimpleBlobLPlugin(t *testing.T) {
	plugin, err := types.NewSimpleBlobLPlugin("test", `root.secret = this.secret.encode("base64")`)
	assert.Nil(t, err)
	data, err := plugin.Transform(map[string]interface{}{
		"secret": "secret",
	})
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"secret": "c2VjcmV0"}, data)
}
