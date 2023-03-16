package plugins_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var bloblPlugin = lo.Must(plugins.NewBloblangPlugin("blobl", `root.secret = this.data.secret.encode("base64")`))

func TestNewSimpleBlobLPlugin(t *testing.T) {
	data, err := bloblPlugin.Execute(context.Background(), secretMessage())

	assert.Nil(t, err)
	assert.Equal(t, secretEncodedMessage(), data)

	plugin, err := plugins.NewBloblangPlugin("test", `
	root.data = this.data
	root.data.test = "test"
	`)
	assert.Nil(t, err)
	data, err = plugin.Execute(context.Background(), complexMessage())
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{
		"Data": float64(1),
		"Next": map[string]any{
			"Data": float64(1),
			"Next": nil,
		},
		"test": "test",
	}, data.Data)
}

func TestNewSimpleBlobLPluginFailureCase(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `some invalid blobl`)
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "failed to parse bloblang template")

	plugin, err = plugins.NewBloblangPlugin("test", `throw("some error")`)
	assert.Nil(t, err)
	_, err = plugin.Execute(context.Background(), emptyMessage())
	assert.ErrorContains(t, err, "some error")
}

func TestNewSimpleBlobLPluginCondition(t *testing.T) {
	plugin, err := plugins.NewBloblangPlugin("test", `this.data.test == "test"`)
	assert.Nil(t, err)
	data, err := plugin.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.True(t, data.Data.(bool))
}
