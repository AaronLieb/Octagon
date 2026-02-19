// Package stream handles the setup and controlling of the tournament stream and the stream sets
package stream

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "stream",
		UsageText: "octagon stream [subcommand]",
		Usage:     "start and automatically manage the stream",
		Aliases:   []string{"c", "conflict"},
		Commands: []*cli.Command{
			WatchCommand(),
			SetupCommand(),
			SelectCommand(),
		},
	}
}
