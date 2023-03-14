package plugins

import (
	"context"
	"errors"
)

type BaseManager[T Plugin] struct {
	plugins map[string]T
}

func NewBaseManager[T Plugin]() Manager[T] {
	return &BaseManager[T]{
		plugins: make(map[string]T),
	}
}

func (m *BaseManager[T]) Add(plugin T) {
	m.plugins[plugin.GetName()] = plugin
}

func (p *BaseManager[T]) Get(name string) (T, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return plugin, errors.New("plugin not found")
	}
	return plugin, nil
}

func (p *BaseManager[T]) Execute(ctx context.Context, name string, data *Message) (*Message, error) {
	plugin, err := p.Get(name)
	if err != nil {
		return nil, err
	}
	return plugin.Execute(ctx, data)
}

type BasePluginManager struct {
	Manager[Plugin]
}

func NewBasePluginManager() *BasePluginManager {
	return &BasePluginManager{
		Manager: NewBaseManager[Plugin](),
	}
}

func (m *BasePluginManager) AddOrchestrator(plugin Plugin) {
	m.Add(&OrchestratorPlugin{
		manager: m,
		plugin:  plugin,
	})
}

type BaseWorkflowManager struct {
	Manager[WorkflowPlugin]
}

func NewBaseWorkflowManager() *BaseWorkflowManager {
	return &BaseWorkflowManager{
		Manager: NewBaseManager[WorkflowPlugin](),
	}
}
