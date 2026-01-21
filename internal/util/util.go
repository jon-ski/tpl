package util

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/jon-ski/tpl/internal/term"
)

// GetInput returns either the file input or stdin if no file flag is given
func GetInput(filepath string) (io.Reader, error) {
	if filepath == "" {
		slog.Debug("file not provided, using stdin")
		return term.GetStdin()
	}
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("failed to open input file", slog.String("err", err.Error()), slog.String("path", filepath))
		return file, fmt.Errorf("failed to open file [%s]: %s", filepath, err)
	}
	return file, nil
}

func getStdin() (io.Reader, error) {
	panic("unimplemented")
}
