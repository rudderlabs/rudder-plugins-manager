package plugins

import (
	"context"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/rudderlabs/rudder-plugins-manager/types"
)

/**
 * This is a base plugin that can be used to create a plugin from a function.
 */
type BasePlugin struct {
	Name     string
	Executor types.Executor
}

func NewBasePlugin(name string, executor types.Executor) *BasePlugin {
	return &BasePlugin{Name: name, Executor: executor}
}

func (p *BasePlugin) Execute(ctx context.Context, data any) (any, error) {
	return p.Executor.Execute(ctx, data)
}

func (p *BasePlugin) GetName() string {
	return p.Name
}

func NewTransformPlugin(name string, fn func(any) (any, error)) types.Plugin {
	return &BasePlugin{
		Name:     name,
		Executor: types.TransformFunc(fn),
	}
}

/**
 * This is transforms the data using bloblang template.
 */
type BloblangPlugin struct {
	Name     string
	executor *bloblang.Executor
}

func NewBloblangPlugin(name string, template string) (*BloblangPlugin, error) {
	executor, err := bloblang.Parse(template)
	if err != nil {
		return nil, err
	}
	return &BloblangPlugin{Name: name, executor: executor}, nil
}

func (p *BloblangPlugin) Execute(_ context.Context, input any) (any, error) {
	return p.executor.Query(input)
}

func (p *BloblangPlugin) GetName() string {
	return p.Name
}
