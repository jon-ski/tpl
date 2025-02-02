package main

import (
	"log/slog"
	"path/filepath"
	"strings"
)

func getFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path)) // Get extension and convert to lower case
	extValue := ext
	switch ext {
	case ".csv":
		extValue = "csv"
	case ".xml":
		extValue = "xml"
	case ".json":
		extValue = "json"
	default:
		return ext
	}
	slog.Debug("file type", slog.String("ext", extValue))
	return extValue
}
