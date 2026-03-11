package cmd

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/jon-ski/cli"
	"github.com/jon-ski/tpl/internal/data"
	"github.com/jon-ski/tpl/internal/template"
	"github.com/jon-ski/tpl/internal/term"
)

var logLevel = new(slog.LevelVar)

// NewRootCmd return the root tpl command
func NewRootCmd() *cli.Command {
	c := cli.NewCommand("", "", "")
	var datasources data.DataSourceList
	var templateParameter string
	c.Flags.Var(&datasources, "d", "named datasource (name=source)")
	c.Flags.StringVar(&templateParameter, "i", "", "-i {{ . }}")

	c.Run = func(ctx *cli.Context, args []string) error {
		// Get datasources, if any
		d, err := data.GetAll(datasources)
		if err != nil {
			return fmt.Errorf("failed to get data sources: %w", err)
		}

		templateString := templateParameter
		// Figure out template source
		if templateParameter == "" {
			// Expect stdin if parameter not provided
			stdin, err := term.GetStdin()
			if err != nil {
				return fmt.Errorf("failed to get stdin: %w", err)
			}
			stdinBytes, err := io.ReadAll(stdin)
			templateString = string(stdinBytes)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		err = template.RunTemplate(templateString, d, ctx.Stdout)
		if err != nil {
			return fmt.Errorf("failed to run template: %w", err)
		}

		return nil
	}

	return c
}

func Execute() {
	cli.Main(NewRootCmd())
}
