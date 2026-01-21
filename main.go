package main

import (
	"log/slog"
	"os"

	"github.com/jon-ski/tpl/cmd"
)

var logLevel = new(slog.LevelVar)

func main() {
	// Set up log
	logLevel.Set(slog.LevelInfo)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cmd.Execute()
}
