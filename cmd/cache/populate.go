package cache

import (
	"context"

	"github.com/AaronLieb/octagon/cache"
	"github.com/urfave/cli/v3"
)

func PopulateCommand() *cli.Command {
	return &cli.Command{
		Name:    "populate",
		Usage:   "Populate player cache from tournament history",
		Aliases: []string{"p"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "Tournament slug (e.g., octagon-136)",
				Value:   "octagon",
			},
			&cli.IntFlag{
				Name:  "history",
				Usage: "Number of tournaments to check",
				Value: 5,
			},
		},
		Action: populateCache,
	}
}

func populateCache(ctx context.Context, cmd *cli.Command) error {
	tournament := cmd.String("tournament")
	history := int(cmd.Int("history"))

	return cache.PopulatePlayerCache(ctx, tournament, history)
}
