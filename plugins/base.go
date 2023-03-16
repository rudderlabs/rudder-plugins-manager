package plugins

import (
	"context"
)

/**
 * This is a base plugin that can be used to create a plugin from a function.
 */
type BasePlugin struct {
	Name     string
	Executor Executor
}

func NewBasePlugin(name string, executor Executor) Plugin {
	return &BasePlugin{Name: name, Executor: executor}
}

func (p *BasePlugin) Execute(ctx context.Context, data *Message) (*Message, error) {
	return p.Executor.Execute(ctx, data)
}

func (p *BasePlugin) GetName() string {
	return p.Name
}

func NewTransformPlugin(name string, fn func(*Message) (*Message, error)) Plugin {
	return &BasePlugin{
		Name:     name,
		Executor: TransformFunc(fn),
	}
}
