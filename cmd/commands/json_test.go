package commands_test

import (
	"testing"

	"github.com/rudderlabs/rudder-transformations/cmd/commands"
	"github.com/rudderlabs/rudder-transformations/utils"
	"github.com/rudderlabs/rudder-transformations/utils/test"
	"github.com/stretchr/testify/assert"
)

type TransformedRecord struct {
	commands.Record
	Default bool `json:"default"`
	Custom  bool `json:"custom"`
}

func TestJSONCmd(t *testing.T) {
	pluginManager := test.GetTestPluginManager()
	command := commands.GetJSONCmd(pluginManager)
	command.SetArgs([]string{"--input", "testdata/input.json", "--output", "generated/output.json", "--provider", "destinations"})
	err := command.Execute()
	assert.Nil(t, err)
	expected, err := utils.ReadJSONFromFile[TransformedRecord]("testdata/expected_output.json")
	assert.Nil(t, err)
	actual, err := utils.ReadJSONFromFile[TransformedRecord]("generated/output.json")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
