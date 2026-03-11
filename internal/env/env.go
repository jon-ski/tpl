package env

import (
	"os"
	"strings"
)

func GetEnv() (map[string]string, error) {
	data := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		data[pair[0]] = pair[1]
	}

	return data, nil
}
