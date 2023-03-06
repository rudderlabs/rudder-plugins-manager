package utils_test

import (
	"errors"
	"testing"

	"github.com/rudderlabs/rudder-plugins-manager/utils"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	v := utils.Must(1, nil)
	assert.Equal(t, 1, v)
	assert.Panics(t, func() { utils.Must[any](nil, errors.New("error")) })
}

func TestStringPtr(t *testing.T) {
	v := utils.StringPtr("test")
	assert.Equal(t, "test", *v)
}
