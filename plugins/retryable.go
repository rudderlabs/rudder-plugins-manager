package plugins

import (
	"context"

	"github.com/cenkalti/backoff/v4"
)

type BaseRetryPolicy struct {
	RetryCount uint64
	Errors     []string
	IsErrorRetryableFunc
}

func (r *BaseRetryPolicy) IsErrorRetryable(err error) bool {
	if r.IsErrorRetryableFunc != nil {
		return r.IsErrorRetryableFunc(err)
	}
	for _, e := range r.Errors {
		if e == err.Error() {
			return true
		}
	}
	return true
}

func (r *BaseRetryPolicy) GetRetryCount() uint64 {
	return r.RetryCount
}

type BaseRetryableExecutor struct {
	Executor
	Policy RetryPolicy
}

func NewBaseRetryableExecutorFromFunc(executor Executor, retryCount uint64, IsErrorRetryable IsErrorRetryableFunc) *BaseRetryableExecutor {
	return &BaseRetryableExecutor{
		Executor: executor,
		Policy: &BaseRetryPolicy{
			RetryCount:           retryCount,
			IsErrorRetryableFunc: IsErrorRetryable,
		},
	}
}

func NewBaseRetryableExecutor(executor Executor, policy RetryPolicy) Executor {
	if policy == nil {
		return executor
	}
	return &BaseRetryableExecutor{
		Executor: executor,
		Policy:   policy,
	}
}

func (e *BaseRetryableExecutor) Execute(ctx context.Context, data *Message) (*Message, error) {
	return backoff.RetryWithData(
		func() (*Message, error) {
			result, err := e.Executor.Execute(ctx, data)
			if err != nil {
				if !e.Policy.IsErrorRetryable(err) {
					err = backoff.Permanent(err)
				}
				return nil, err
			}
			return result, nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), e.Policy.GetRetryCount()))
}

func NewBaseRetryablePlugin(plugin Plugin, policy RetryPolicy) Plugin {
	return &BasePlugin{
		Name:     plugin.GetName(),
		Executor: NewBaseRetryableExecutor(plugin, policy),
	}
}

func NewBasePluginWithRetryPolicy(name string, executor Executor, policy RetryPolicy) Plugin {
	return &BasePlugin{
		Name:     name,
		Executor: NewBaseRetryableExecutor(executor, policy),
	}
}
