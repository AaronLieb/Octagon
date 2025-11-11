/* Package conflicts provides commands for creating, viewing, and managing manual conflicts */
package conflicts

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "conflicts",
		UsageText: "octagon conflicts [subcommand]",
		Usage:     "View and manage conflicts",
		Aliases:   []string{"c", "conflict"},
		Commands: []*cli.Command{
			createCommand(),
			listCommand(),
		},
	}
}
