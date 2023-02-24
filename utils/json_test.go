package utils_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rudderlabs/rudder-transformations/utils"
	"github.com/stretchr/testify/assert"
)

type testRecord struct {
	Test string `json:"test"`
}

func TestReadJSONFromFileHappyCase(t *testing.T) {
	records := []testRecord{
		{
			Test: "test1",
		},
		{
			Test: "test2",
		},
		{
			Test: "test3",
		},
	}
	data, err := json.Marshal(records)
	assert.Nil(t, err)
	err = os.WriteFile("generated.json", data, 0o644)
	assert.Nil(t, err)
	actual, err := utils.ReadRecordsFromJSONFile[testRecord]("generated.json")
	assert.Nil(t, err)
	assert.Equal(t, records, actual)
}

func TestReadJSONFromFileErrorCases(t *testing.T) {
	_, err := utils.ReadRecordsFromJSONFile[testRecord]("non-existing-file.json")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "no such file or directory")

	_, err = utils.ReadRecordsFromJSONFile[testRecord]("../testdata/input.txt")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid character")
}
