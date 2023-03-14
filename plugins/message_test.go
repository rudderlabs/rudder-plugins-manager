package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageClone(t *testing.T) {
	testMsg := testMessage()
	clone := testMsg.Clone()
	assert.Equal(t, testMsg, clone)
	clone.SetMetadata("test", "test")
	assert.NotEqual(t, testMsg, clone)
}

func TestMessageMetadata(t *testing.T) {
	testMsg := testMessage()
	testMsg.SetMetadata("test", "test")
	value, ok := testMsg.GetMetadata("test")
	assert.True(t, ok)
	assert.Equal(t, "test", value)
}

func TestMessageGetBool(t *testing.T) {
	testMsg := testMessage()
	testMsg.Data = true
	value, err := testMsg.GetBool()
	assert.Nil(t, err)
	assert.True(t, value)

	testMsg.Data = "test"
	value, err = testMsg.GetBool()
	assert.NotNil(t, err)
	assert.False(t, value)
}
