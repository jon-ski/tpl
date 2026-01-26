package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jon-ski/cli"
	"github.com/jon-ski/tpl/internal/csv"
	"github.com/jon-ski/tpl/internal/template"
	"github.com/jon-ski/tpl/internal/util"
)

// CmdCsv is the cmd subcommand used to ingest csv data before running the template
//
// tpl csv '{{ .csv.rows }}' input.csv
// cat input.csv | tpl csv '{{ .csv.rows }}'
// cat input.csv | tpl csv -t main.tmpl
// However, initial function will expect stdin rather than file and first argument always the template
type CmdCsv struct {
	Cmd *cli.Command

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

	if len(args) < 1 {
		return fmt.Errorf("expected template argument")
	}
	templateString := args[0]

	// TODO: get file path from flag if provided
	input, err := util.GetInput("")
	if err != nil {
		return fmt.Errorf("failed to setup input data: %w", err)
	}

	slog.Debug("creating data map")
	data, err := csv.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to create data map: %w", err)
	}

	slog.Debug("running template")
	err = template.RunTemplate(templateString, data, ctx.Stdout)
	if err != nil {
		return fmt.Errorf("failed to run template: %w", err)
	}

	return nil
}
