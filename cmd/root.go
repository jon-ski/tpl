package cmd

import (
	"log/slog"

	"github.com/jon-ski/cli"
)

var logLevel = new(slog.LevelVar)

func NewRootCmd() *cli.Command {
	c := cli.NewCommand("", "", "")
	c.Add(NewCmdCsv())
	c.Add(NewCmdJson())
	c.Add(NewCmdVersion())
	return c
}

func Execute() {
	cli.Main(NewRootCmd())
}
