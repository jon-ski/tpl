package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	icsv "github.com/jon-ski/tpl/internal/csv"
	ijson "github.com/jon-ski/tpl/internal/json"
	ixml "github.com/jon-ski/tpl/internal/xml"
)

type Flags struct {
	inputFile    *string
	templateFile *string

	debug *bool
}

func flagParse() *Flags {
	flags := &Flags{
		inputFile:    flag.String("i", "", "input file path"),
		templateFile: flag.String("t", "main.tmpl", "template file path (default: 'main.tmpl')"),
		debug:        flag.Bool("d", false, "sets log level to debug"),
	}
	flag.Parse()
	return flags
}

var logLevel = new(slog.LevelVar)

func main() {
	// Parse flags
	flags := flagParse()

	// Set up log
	logLevel.Set(slog.LevelInfo)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cmd.Execute()
}
