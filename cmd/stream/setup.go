package stream

import (
	"context"

	"github.com/urfave/cli/v3"
)

func SetupCommand() *cli.Command {
	return &cli.Command{
		Name:      "watch",
		Aliases:   []string{"w"},
		Usage:     "sets up tsh and obs",
		UsageText: "octagon stream setup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "Tournament slug",
				Value:   "octagon",
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage:   "Event name",
				Value:   "ultimate-singles",
			},
		},
		Action: setup,
	}
}

func setup(ctx context.Context, cmd *cli.Command) error {
	return nil
}
