package plugins

import (
	"github.com/samber/lo"
	"github.com/vmihailenco/msgpack/v5"
)

func Clone[T any](value T) (T, error) {
	var newValue T
	valueBytes, err := msgpack.Marshal(value)
	if err != nil {
		return newValue, err
	}
	lo.Must0(msgpack.Unmarshal(valueBytes, &newValue))
	return newValue, nil
}
