package plugins

import (
	"context"
)

/**
 * This file contains the interface that all plugins must implement.
 * The interface is used by the plugin manager to load and use plugins.
 */

type Executor interface {
	Execute(context.Context, *Message) (*Message, error)
}

type Plugin interface {
	Executor
	GetName() string
}

type IsErrorRetryableFunc func(error) bool

type RetryableExecutor interface {
	Executor
	GetRetryCount() int
	IsErrorRetryable(error) bool
}

var AllErrorsRetryable IsErrorRetryableFunc = func(error) bool { return true }

type ExecuteFunc func(context.Context, *Message) (*Message, error)

func (f ExecuteFunc) Execute(ctx context.Context, data *Message) (*Message, error) {
	return f(ctx, data)
}

type TransformFunc func(*Message) (*Message, error)

func (f TransformFunc) Execute(_ context.Context, data *Message) (*Message, error) {
	return f(data)
}

type Manager[T Plugin] interface {
	ExecutionManager
	Get(name string) (T, error)
	Has(name string) bool
	Add(plugin T)
}

type ExecutionManager interface {
	Execute(ctx context.Context, name string, data *Message) (*Message, error)
}

type StepType string

const (
	BloblangStep   StepType = "bloblang"
	PluginStep     StepType = "plugin"
	ExpressionStep StepType = "expression"
	UnknownStep    StepType = "unknown"
)

type RetryPolicy interface {
	IsErrorRetryable(error) bool
	GetRetryCount() uint64
}

type StepPlugin interface {
	Plugin
	GetType() StepType
	ShouldExecute(*Message) (bool, error)
	ShouldReturn() bool
	ShouldContinue() bool
	GetRetryPolicy() (RetryPolicy, bool)
}

type WorkflowPlugin interface {
	Plugin
	GetVersion() int
	GetSteps() []StepPlugin
	GetStep(name string) (StepPlugin, error)
	ExecuteStep(ctx context.Context, stepName string, data *Message) (*Message, error)
}

type WorkflowConfig struct {
	Name    string       `json:"name" yaml:"name"`
	Version int          `json:"version" yaml:"version"`
	Steps   []StepConfig `json:"steps" yaml:"steps"`
}

func (c *WorkflowConfig) GetVersion() int {
	if c.Version == 0 {
		return 1
	}
	return c.Version
}

type StepConfig struct {
	Name       string           `json:"name" yaml:"name"`
	CheckBlobl string           `json:"check_blobl" yaml:"check_blobl"`
	CheckExpr  string           `json:"check_expr" yaml:"check_expr"`
	Return     bool             `json:"return" yaml:"return"`
	Continue   bool             `json:"continue" yaml:"continue"`
	Plugin     string           `json:"plugin" yaml:"plugin"`
	Bloblang   string           `json:"bloblang" yaml:"bloblang"`
	Expr       string           `json:"expr" yaml:"expr"`
	Retry      *BaseRetryPolicy `json:"retry" yaml:"retry"`
}

func (c *StepConfig) GetType() StepType {
	if c.Bloblang != "" {
		return BloblangStep
	} else if c.Plugin != "" {
		return PluginStep
	} else if c.Expr != "" {
		return ExpressionStep
	}
	return UnknownStep
}

type (
	PluginManager   Manager[Plugin]
	WorkflowManager Manager[WorkflowPlugin]
)
