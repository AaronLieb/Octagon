package attendees

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "attendees",
		UsageText: "octagon attendees [subcommand]",
		Usage:     "Manage tournament attendees",
		Aliases:   []string{"a"},
		Commands: []*cli.Command{
			addCommand(),
			listCommand(),
			redemptionInfoCommand(),
		},
	}
}
