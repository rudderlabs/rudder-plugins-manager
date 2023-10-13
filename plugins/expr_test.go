package plugins_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var exprPlugin = lo.Must(plugins.NewExpressionPlugin("expr", `{secret: toBase64(data.secret)}`))

func TestNewSimpleExprPlugin(t *testing.T) {
	assert.Equal(t, "expr", exprPlugin.GetName())
	data, err := exprPlugin.Execute(context.Background(), secretMessage())

	assert.Nil(t, err)
	assert.Equal(t, secretEncodedMessage(), data)
}

func TestNewSimpleExprPluginFailureCase(t *testing.T) {
	plugin, err := plugins.NewExpressionPlugin("test", `some invalid expr`)
	assert.NotNil(t, err)
	assert.Nil(t, plugin)
	assert.ErrorContains(t, err, "failed to parse expression")

	plugin, err = plugins.NewExpressionPlugin("test", `data.non_existants.test`)
	assert.Nil(t, err)
	_, err = plugin.Execute(context.Background(), emptyMessage())
	assert.ErrorContains(t, err, "cannot fetch test")
}

func TestNewSimpleExprPluginCondition(t *testing.T) {
	plugin, err := plugins.NewExpressionPlugin("test", `data.test == "test"`)
	assert.Nil(t, err)
	data, err := plugin.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.True(t, data.Data.(bool))
}
