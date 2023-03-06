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
	GetPlugin(name string) (Plugin, error)
}
