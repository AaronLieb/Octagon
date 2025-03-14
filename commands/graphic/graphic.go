package seed

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "graphic",
		Usage:   "Generate a top 8 graphic",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
			&cli.BoolFlag{
				Name:    "redemption",
				Aliases: []string{"r"},
				Usage:   "Whether you are seeding for redemption bracket or not",
			},
		},
		Action: graphic,
	}
}

func graphic(ctx context.Context, cmd *cli.Command) error {
	log.Debug("graphic")

	return nil
}
