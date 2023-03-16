package plugins

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/huandu/go-clone"
	"github.com/samber/lo"
)

type Message struct {
	Data any `json:"data"`
	// This will be used in workflows to store the original input message
	Input    any            `json:"input"`
	Metadata map[string]any `json:"metadata"`
}

func NewMessage(data any) *Message {
	msg := Message{
		Data:     data,
		Input:    data,
		Metadata: make(map[string]any),
	}
	return msg.Clone()
}

func (m *Message) Clone() *Message {
	return clone.Slowly(m).(*Message)
}

func (m *Message) WithMetadata(key string, value any) *Message {
	m.Metadata[key] = value
	return m
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
	var result map[string]interface{}
	lo.Must0(json.Unmarshal(lo.Must(json.Marshal(m)), &result))
	return result
}
