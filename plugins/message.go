package plugins

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/huandu/go-clone"
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
	mBytes, _ := json.Marshal(m)
	json.Unmarshal(mBytes, &result)
	return result
}
