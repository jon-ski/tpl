package cmd

import (
	"fmt"

	"github.com/jon-ski/cli"
	"github.com/jon-ski/tpl/internal/version"
)

func NewCmdVersion() *cli.Command {
	c := cli.NewCommand("version", "print version information", "")
	c.Run = func(ctx *cli.Context, args []string) error {
		fmt.Printf("Version=%s\nCommit=%s\nDate=%s\n", version.Version, version.Commit, version.Date)
		return nil
	}
	return c
}
