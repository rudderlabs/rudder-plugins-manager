package plugins

import (
	"fmt"

	"github.com/huandu/go-clone"
)

type Message struct {
	Data any `mapstructure:"data"`
	// This will be used in workflows to store the original input message
	Input    any            `mapstructure:"input"`
	Metadata map[string]any `mapstructure:"metadata"`
}

func NewMessage(data any) *Message {
	return &Message{
		Data:     data,
		Input:    data,
		Metadata: make(map[string]any),
	}
}

func (m *Message) Clone() *Message {
	return clone.Clone(m).(*Message)
}

func (m *Message) SetMetadata(key string, value any) {
	m.Metadata[key] = value
}

func (m *Message) GetMetadata(key string) (any, bool) {
	value, ok := m.Metadata[key]
	return value, ok
}

func (m *Message) GetBool() (bool, error) {
	value, ok := m.Data.(bool)
	if !ok {
		return false, fmt.Errorf("data is not bool")
	}
	return value, nil
}

func (m *Message) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data":     m.Data,
		"input":    m.Input,
		"metadata": m.Metadata,
	}
}
