package term

import (
	"errors"
	"io"
	"os"
)

func GetStdin() (io.Reader, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("stdin is a terminal; pipe or redirect input")
	}
	return os.Stdin, nil
}
