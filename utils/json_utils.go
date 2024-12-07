package utils

import (
	"encoding/json"
	"errors"
)

func MarshalObject(input interface{}) ([]byte, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
