package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jon-ski/cli"
	"github.com/jon-ski/tpl/internal/json"
	"github.com/jon-ski/tpl/internal/template"
	"github.com/jon-ski/tpl/internal/util"
)

type CmdJson struct {
	Cmd *cli.Command

	TemplatePath string
	Verbose      bool
}

func NewCmdJson() *cli.Command {
	c := &CmdJson{}
	c.Cmd = cli.NewCommand("json", "ingest json data", "[input file]")
	c.Cmd.Flags.StringVar(&c.TemplatePath, "t", "main.tmpl", "template root file")
	c.Cmd.Flags.BoolVar(&c.Verbose, "v", false, "verbose output")
	c.Cmd.Run = c.run
	return c.Cmd
}

func (c *CmdJson) run(ctx *cli.Context, args []string) error {
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
	data, err := json.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to parse json input: %w", err)
	}

	slog.Debug("running template")
	err = template.RunTemplate(c.TemplatePath, data)
	if err != nil {
		return fmt.Errorf("failed to run template: %w", err)
	}

	return nil
}
