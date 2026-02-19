package bracket

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/seeding"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const (
	desc = `Will fetch player ratings from an external rating database
and seed the players accordingly. It will then attempt to read all
player conflicts and generate a variation of the original seeding
that minimizes seeding changes while maximizing conflict resolution`
)

func seedCommand() *cli.Command {
	return &cli.Command{
		Name:        "seed",
		Usage:       "Seeds an event",
		Description: desc,
		Aliases:     []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage:   "The event name. Example: 'ultimate-singles'",
			},
			&cli.BoolFlag{
				Name:    "redemption",
				Aliases: []string{"r"},
				Usage:   "Whether you are seeding for redemption bracket or not",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "A conflict file to read",
			},
		},
		Action: seed,
	}
}

func seed(ctx context.Context, cmd *cli.Command) error {
	log.Debug("seed")

	tournamentShortSlug := cmd.String("tournament")
	redemption := cmd.Bool("redemption")

	var conflictFiles []string
	if cmd.String("file") != "" {
		conflictFiles = append(conflictFiles, cmd.String("file"))
	}

	players, err := seeding.GenerateSeeding(ctx, tournamentShortSlug, redemption, conflictFiles)
	if err != nil {
		return err
	}

	var input string
	fmt.Println("Publish seeding? (y/N)")
	_, err = fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}

	if input != "y" && input != "Y" {
		log.Info("Cancelling seeding...")
		return nil
	}

	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	event := startgg.EventUltimateSingles
	if redemption {
		event = startgg.EventRedemptionBracket
	}

	if len(cmd.String("event")) > 0 {
		event = cmd.String("event")
	}

	slug := fmt.Sprintf(startgg.EventSlugFormat, tournamentSlug, event)
	err = seeding.PublishSeeding(ctx, slug, players)
	if err != nil {
		return err
	}

	log.Infof("Successfully seeded %d players for %s", len(players), slug)
	return nil
}
