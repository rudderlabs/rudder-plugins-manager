package types_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/types"
	"github.com/rudderlabs/rudder-plugins-manager/utils"
	"github.com/stretchr/testify/assert"
)

func TestExecuteFunc(t *testing.T) {
	executor := types.ExecuteFunc(func(ctx context.Context, data any) (any, error) {
		return data, nil
	})
	assert.NotNil(t, executor)
	assert.Equal(t, "test", utils.Must(executor.Execute(context.Background(), "test")))
}

func TestTransformFunc(t *testing.T) {
	executor := types.TransformFunc(func(data any) (any, error) {
		return data, nil
	})
	assert.NotNil(t, executor)
	assert.Equal(t, "test", utils.Must(executor.Execute(context.Background(), "test")))
}
