package plugins

import (
	"context"
	"fmt"
)

type BaseManager[T Plugin] struct {
	plugins map[string]T
	Type    string
}

func NewBaseManager[T Plugin](typeVal string) Manager[T] {
	return &BaseManager[T]{
		Type:    typeVal,
		plugins: make(map[string]T),
	}
}

func (m *BaseManager[T]) Add(plugin T) {
	m.plugins[plugin.GetName()] = plugin
}

func (m *BaseManager[T]) Has(name string) bool {
	_, ok := m.plugins[name]
	return ok
}

func (m *BaseManager[T]) Get(name string) (T, error) {
	plugin, ok := m.plugins[name]
	if !ok {
		return plugin, fmt.Errorf("%s %s not found", name, m.Type)
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
		Manager: NewBaseManager[Plugin]("plugin"),
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
		Manager: NewBaseManager[WorkflowPlugin]("workflow"),
	}
}
