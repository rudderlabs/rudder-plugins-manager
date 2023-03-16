package plugins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/plugins"
	"github.com/stretchr/testify/assert"
)

type failingExecutor struct {
	message          string
	passAfterRetries int
}

func newFailingExecutor(message string, passAfterRetries int) *failingExecutor {
	return &failingExecutor{message: message, passAfterRetries: passAfterRetries}
}

func (f *failingExecutor) Execute(ctx context.Context, msg *plugins.Message) (*plugins.Message, error) {
	if f.passAfterRetries > 0 {
		f.passAfterRetries--
		return nil, errors.New(f.message)
	}
	return msg, nil
}

func retryableError(err error) bool {
	return err.Error() == "retryable"
}

func newFailingPlugin(name, message string, passAfterRetries int) plugins.Plugin {
	return plugins.NewBasePlugin(name, newFailingExecutor(message, passAfterRetries))
}

func TestNewBaseRetryableExecutor(t *testing.T) {
	retryPolicy := &plugins.BaseRetryPolicy{
		RetryCount:           1,
		IsErrorRetryableFunc: plugins.AllErrorsRetryable,
	}
	// Retry will succeed
	executor := failingExecutor{message: "failed", passAfterRetries: 1}
	retryableExecutor := plugins.NewBaseRetryableExecutor(&executor, retryPolicy)
	data, err := retryableExecutor.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)

	// Retry will succeed as we will try for all errors
	// when Errors or IsErrorRetryableFunc is nil
	executor = failingExecutor{message: "failed", passAfterRetries: 1}
	retryableExecutor = plugins.NewBaseRetryableExecutor(&executor, &plugins.BaseRetryPolicy{RetryCount: 2})
	data, err = retryableExecutor.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)

	// Retry will fail because of nil policy
	executor = failingExecutor{message: "failed", passAfterRetries: 1}
	retryableExecutor = plugins.NewBaseRetryableExecutor(&executor, nil)
	_, err = retryableExecutor.Execute(context.Background(), testMessage())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed")

	// Retry will fail due to max retries
	executor = failingExecutor{message: "failed", passAfterRetries: 5}
	retryableExecutor = plugins.NewBaseRetryableExecutor(&executor, retryPolicy)
	_, err = retryableExecutor.Execute(context.Background(), testMessage())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed")

	// Now retry will succeed because we are throwing a retryable error
	retryPolicy = &plugins.BaseRetryPolicy{
		RetryCount: 2,
		Errors:     []string{"retryable"},
	}
	executor = failingExecutor{message: "retryable", passAfterRetries: 1}
	retryableExecutor = plugins.NewBaseRetryableExecutor(&executor, retryPolicy)
	data, err = retryableExecutor.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestNewBaseRetryableExecutorFromFunc(t *testing.T) {
	// Now retry will fail because we are throwing a non-retryable error
	executor := failingExecutor{message: "failed", passAfterRetries: 1}
	retryableExecutor := plugins.NewBaseRetryableExecutorFromFunc(&executor, 1, retryableError)
	_, err := retryableExecutor.Execute(context.Background(), testMessage())
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed")

	// Now retry will succeed because we are throwing a retryable error
	executor = failingExecutor{message: "retryable", passAfterRetries: 1}
	retryableExecutor = plugins.NewBaseRetryableExecutorFromFunc(&executor, 1, retryableError)
	data, err := retryableExecutor.Execute(context.Background(), testMessage())
	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestNewBaseRetryablePlugin(t *testing.T) {
	// Retry will succeed
	retryablePlugin := plugins.NewBaseRetryablePlugin(
		newFailingPlugin("test", "failed", 1),
		&plugins.BaseRetryPolicy{RetryCount: 1},
	)
	assert.Equal(t, "test", retryablePlugin.GetName())
	data, err := retryablePlugin.Execute(context.Background(), testMessage())

	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}

func TestNewBasePluginWithRetryPolicy(t *testing.T) {
	// Retry will succeed
	retryablePlugin := plugins.NewBasePluginWithRetryPolicy(
		"test",
		newFailingExecutor("failed", 1),
		&plugins.BaseRetryPolicy{RetryCount: 1},
	)
	assert.Equal(t, "test", retryablePlugin.GetName())
	data, err := retryablePlugin.Execute(context.Background(), testMessage())

	assert.Nil(t, err)
	assert.Equal(t, testMessage(), data)
}
