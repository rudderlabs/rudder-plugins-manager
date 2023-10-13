package plugins

import (
	"context"
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

/**
 * This is transforms the data using bloblang template.
 */
type ExpressionPlugin struct {
	Name     string
	program *vm.Program
}

func NewExpressionPlugin(name, template string) (*ExpressionPlugin, error) {
	program, err := expr.Compile(template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}
	return &ExpressionPlugin{Name: name, program: program}, nil
}

func (p *ExpressionPlugin) Execute(_ context.Context, input *Message) (*Message, error) {
	inputMap := input.ToMap()
	log.Debug().Any("input", inputMap).Msg("Executing Expression plugin")
	data, err := expr.Run(p.program, inputMap)
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

func (p *ExpressionPlugin) GetName() string {
	return p.Name
}
