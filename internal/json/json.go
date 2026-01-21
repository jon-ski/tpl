package json

import (
	"encoding/json"
	"fmt"
	"io"
)

func Parse(input io.Reader) (any, error) {
	b, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input data: %w", err)
	}
	data := make(map[string]any)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %w", err)
	}
	return data, nil
}
