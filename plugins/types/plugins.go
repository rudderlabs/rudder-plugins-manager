package types

import (
	"github.com/benthosdev/benthos/v4/public/bloblang"
)

/**
 * This is a simple plugin that can be used to create a plugin from a function.
 */
type SimplePlugin struct {
	Name        string
	Transformer Transformer
}

func NewSimplePlugin(name string, transformFunc func(any) (any, error)) *SimplePlugin {
	return &SimplePlugin{Name: name, Transformer: TransformerFunc(transformFunc)}
}

func (p *SimplePlugin) GetTransformer(any) (Transformer, error) {
	return p.Transformer, nil
}

func (p *SimplePlugin) Transform(data any) (any, error) {
	return p.Transformer.Transform(data)
}

func (p *SimplePlugin) GetName() string {
	return p.Name
}

/**
 * This is transforms the data using bloblang template.
 */
type BlobLTransformer struct {
	executor *bloblang.Executor
}

func NewBlobLTransformer(template string) (*BlobLTransformer, error) {
	executor, err := bloblang.Parse(template)
	if err != nil {
		return nil, err
	}
	return &BlobLTransformer{executor: executor}, nil
}

func (t *BlobLTransformer) Transform(data any) (any, error) {
	return t.executor.Query(data)
}

func NewSimpleBlobLPlugin(name, template string) (*SimplePlugin, error) {
	transformer, err := NewBlobLTransformer(template)
	if err != nil {
		return nil, err
	}
	return &SimplePlugin{Name: name, Transformer: transformer}, nil
}
