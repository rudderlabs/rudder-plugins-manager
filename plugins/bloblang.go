package plugins

import (
	"context"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

/**
 * This is transforms the data using bloblang template.
 */
type BloblangPlugin struct {
	Name     string
	executor *bloblang.Executor
}

func NewBloblangPlugin(name, template string) (*BloblangPlugin, error) {
	executor, err := bloblang.Parse(template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bloblang template: %w", err)
	}
	return &BloblangPlugin{Name: name, executor: executor}, nil
}

func (p *BloblangPlugin) Execute(_ context.Context, input *Message) (*Message, error) {
	inputMap := input.ToMap()
	log.Debug().Any("input", inputMap).Msg("Executing bloblang plugin")
	data, err := p.executor.Query(inputMap)
	if err != nil {
		return nil, err
	}
	var result Message
	var dataBytes []byte
	if dataBytes, err = json.Marshal(data); err == nil {
		err = json.Unmarshal(dataBytes, &result)
	}
	if err != nil || (result.Metadata == nil && result.Data == nil) {
		return NewMessage(data), nil
	}
	return &result, nil
}

func (p *BloblangPlugin) GetName() string {
	return p.Name
}
