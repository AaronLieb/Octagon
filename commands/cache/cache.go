package cache

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "cache",
		Usage:   "cache commands ",
		Aliases: []string{"c"},
		Commands: []*cli.Command{
			ClearCommand(),
		},
	}
}
