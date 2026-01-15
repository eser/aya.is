package caching

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidType = errors.New("invalid type")

//nolint:ireturn
func FromBytes[T any](bytes []byte) (T, error) {
	var result T

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return result, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}

	return result, nil
}

func ToBytes[T any](value T) ([]byte, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}

	return bytes, nil
}
