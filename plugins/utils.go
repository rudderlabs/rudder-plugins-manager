package plugins

import (
	"bytes"
	"encoding/gob"

	"github.com/samber/lo"
)

func init() {
	gob.Register(map[string]interface{}{})
}

func Clone[T any](value T) (T, error) {
	var newValue T
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer) // Will write to buffer.
	dec := gob.NewDecoder(&buffer) // Will read from buffer.
	if err := enc.Encode(value); err != nil {
		return newValue, err
	}
	lo.Must0(dec.Decode(&newValue))
	return newValue, nil
}
