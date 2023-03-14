package plugins

import (
	"context"
)

/**
 * This file contains the interface that all plugins must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */

type Plugin interface {
	GetName() string
	Execute(context.Context, *Message) (*Message, error)
}

type Executor interface {
	Execute(context.Context, *Message) (*Message, error)
}

type ExecuteFunc func(context.Context, *Message) (*Message, error)

func (f ExecuteFunc) Execute(ctx context.Context, data *Message) (*Message, error) {
	return f(ctx, data)
}

type TransformFunc func(*Message) (*Message, error)

func (f TransformFunc) Execute(_ context.Context, data *Message) (*Message, error) {
	return f(data)
}

type Manager[T Plugin] interface {
	Get(name string) (T, error)
	Add(plugin T)
	Execute(ctx context.Context, name string, data *Message) (*Message, error)
}

type StepType string

const (
	BloblangStep StepType = "bloblang"
	PluginStep   StepType = "plugin"
	UnknownStep  StepType = "unknown"
)

type StepPlugin interface {
	Plugin
	GetType() StepType
	ShouldExecute(*Message) (bool, error)
	ShouldReturn() bool
	ShouldContinue() bool
}

type WorkflowPlugin interface {
	Plugin
	GetSteps() []StepPlugin
	GetStep(name string) (StepPlugin, error)
	ExecuteStep(ctx context.Context, stepName string, data *Message) (*Message, error)
}

type WorkflowConfig struct {
	Name  string       `json:"name" yaml:"name"`
	Steps []StepConfig `json:"steps" yaml:"steps"`
}

type StepConfig struct {
	Name     string `json:"name" yaml:"name"`
	Check    string `json:"check" yaml:"check"`
	Return   bool   `json:"return" yaml:"return"`
	Continue bool   `json:"continue" yaml:"continue"`
	Plugin   string `json:"plugin" yaml:"plugin"`
	Bloblang string `json:"bloblang" yaml:"bloblang"`
}

func (c *StepConfig) GetType() StepType {
	if c.Bloblang != "" {
		return BloblangStep
	} else if c.Plugin != "" {
		return PluginStep
	}
	return UnknownStep
}

type PluginManager Manager[Plugin]
type WorkflowManager Manager[WorkflowPlugin]
