package attendees

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "attendees",
		Usage:   "View and modify tournament/event attendees ",
		Aliases: []string{"a"},
		Commands: []*cli.Command{
			AddCommand(),
			ListCommand(),
			RedemptionInfoCommand(),
		},
	}
}
