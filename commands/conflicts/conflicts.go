package conflicts

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "conflicts",
		Description: "manage conflicts",
		Aliases:     []string{"c"},
		Commands: []*cli.Command{
			addCommand(),
			listCommand(),
		},
	}
}
