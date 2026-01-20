package csv

import (
	"encoding/csv"
	"fmt"
	"io"
)

func Parse(input io.Reader) (interface{}, error) {
	// Setup reader
	reader := csv.NewReader(input)

	// Read all (to get line count)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv lines: %w", err)
	}

	// Setup data map
	data := make([]map[string]interface{}, len(rows[1:]))

	// Setup header
	header := rows[0]
	for i := range header {
		if header[i] == "" {
			header[i] = fmt.Sprintf("col-%02d", err)
		}
	}

	// Loop over rows, filling out the data map
	for i, row := range rows[1:] {
		rowData := make(map[string]interface{})
		for j, value := range row {
			rowData[header[j]] = value
		}
		data[i] = rowData
	}
	return data, nil
}
