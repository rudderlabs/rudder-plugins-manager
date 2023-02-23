package utils

import (
	"encoding/json"
	"os"
)

func ReadJSONFromFile[T any](jsonFile string) ([]T, error) {
	var records []T
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}
