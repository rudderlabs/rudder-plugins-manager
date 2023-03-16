package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/stretchr/testify/assert"
)

var testPlugin = plugins.NewTransformPlugin("test", func(msg *plugins.Message) (*plugins.Message, error) {
	dataMap, ok := msg.Data.(map[string]any)
	if !ok {
		return nil, errors.New("data is not a map")
	}
	dataMap["test"] = "test"
	return plugins.NewMessage(dataMap), nil
})

func emptyMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{})
}

func testMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"test": "test"})
}

func complexMessage() *plugins.Message {
	type linkedList struct {
		Data int
		Next *linkedList
	}
	return plugins.NewMessage(linkedList{
		Data: 1,
		Next: &linkedList{
			Data: 1,
		},
	})
}

func secretMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"secret": "secret"})
}

func secretEncodedMessage() *plugins.Message {
	return plugins.NewMessage(map[string]any{"secret": "c2VjcmV0"})
}

var badPlugin = plugins.NewTransformPlugin("bad", func(msg *plugins.Message) (*plugins.Message, error) {
	return nil, errors.New("bad plugin")
})

func TestNewSimplePlugin(t *testing.T) {
	data, err := testPlugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestNewBasePlugin(t *testing.T) {
	plugin := plugins.NewBasePlugin("test",
		plugins.ExecuteFunc(
			func(ctx context.Context, data *plugins.Message) (*plugins.Message, error) {
				return data, nil
			},
		),
	)
	assert.NotNil(t, plugin)
	data, err := plugin.Execute(context.Background(), emptyMessage())
	assert.Nil(t, err)
	assert.Equal(t, emptyMessage(), data)
	assert.Equal(t, "test", plugin.GetName())
}
