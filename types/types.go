package types

import (
	"context"
)

/**
 * This file contains the interface that all plugins must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */
type Plugin interface {
	GetName() string
	Execute(context.Context, any) (any, error)
}

type Executor interface {
	Execute(context.Context, any) (any, error)
}

type ExecuteFunc func(context.Context, any) (any, error)

func (f ExecuteFunc) Execute(ctx context.Context, data any) (any, error) {
	return f(ctx, data)
}

type TransformFunc func(any) (any, error)

func (f TransformFunc) Execute(_ context.Context, data any) (any, error) {
	return f(data)
}

type PluginManager interface {
	Get(name string) (Plugin, error)
	Add(plugin Plugin)
	Execute(ctx context.Context, name string, data any) (any, error)
}

type Pipeline interface {
	GetName() string
	Start(ctx context.Context) error
	Submit(ctx context.Context, data any) error
}

type PipelineManager interface {
	Get(name string) (Pipeline, error)
	Add(pipeline Pipeline)
	Start(ctx context.Context) error
	Submit(ctx context.Context, name string, data any) error
}
