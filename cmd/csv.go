package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jon-ski/cli"
	"github.com/jon-ski/tpl/internal/csv"
	"github.com/jon-ski/tpl/internal/template"
	"github.com/jon-ski/tpl/internal/util"
)

type CmdCsv struct {
	Cmd *cli.Command

	InputPath    string
	TemplatePath string
	Verbose      bool
}

func NewCmdCsv() *cli.Command {
	c := &CmdCsv{}
	c.Cmd = cli.NewCommand("csv", "ingest csv data", "[input file]")
	c.Cmd.Flags.StringVar(&c.TemplatePath, "t", "main.tmpl", "template root file path")
	c.Cmd.Flags.BoolVar(&c.Verbose, "v", false, "verbose output")
	c.Cmd.Run = c.run
	return c.Cmd
}

func (c *CmdCsv) run(ctx *cli.Context, args []string) error {
	if c.Verbose {
		logLevel.Set(slog.LevelDebug)
	}

	if c.TemplatePath == "" {
		return fmt.Errorf("no template file specified")
	}

	inputPath := ""
	if len(args) > 0 {
		inputPath = args[0]
	}

	slog.Debug("initializing input reader")
	input, err := util.GetInput(inputPath)
	if err != nil {
		return fmt.Errorf("failed to setup input data: %w", err)
	}

	slog.Debug("creating data map")
	data, err := csv.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to create data map: %w", err)
	}

	slog.Debug("running template")
	err = template.RunTemplate(c.TemplatePath, data)
	if err != nil {
		return fmt.Errorf("failed to run template: %w", err)
	}

	return nil
}
