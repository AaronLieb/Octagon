/* Package cache provides commands for viewing and manipulating the cache */
package cache

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "cache",
		UsageText: "octagon cache [subcommand]",
		Usage:     "Manage the internal cache",
		Commands: []*cli.Command{
			ClearCommand(),
			PopulateCommand(),
			ListCommand(),
		},
	}
}
