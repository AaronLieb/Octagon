package bracketcmd

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "bracket",
		UsageText: "octagon bracket [subcommand]",
		Usage:     "Manage tournament bracket",
		Aliases:   []string{"b"},
		Commands: []*cli.Command{
			printCommand(),
			seedCommand(),
		},
	}
}
