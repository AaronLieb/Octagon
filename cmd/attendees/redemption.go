package attendees

import (
	"context"
	"fmt"
	"slices"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const DefaultCutoff = 16

func redemptionInfoCommand() *cli.Command {
	return &cli.Command{
		Name:    "redemption",
		Usage:   "Lists players eliminated from main bracket but not signed up for redemption",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
			&cli.IntFlag{
				Name:    "cutoff",
				Aliases: []string{"c"},
				Usage:   "Placement cutoff for redemption eligibility",
				Value:   DefaultCutoff,
			},
		},
		Action: RedemptionInfo,
	}
}

func RedemptionInfo(ctx context.Context, cmd *cli.Command) error {
	log.Debug("Redemption Info")

	tournamentShortSlug := cmd.String("tournament")
	cutoff := cmd.Int("cutoff")

	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	redEventSlug := fmt.Sprintf(startgg.EventSlugFormat, tournamentSlug, startgg.EventRedemptionBracket)
	redEntrantsResp, err := startgg.GetEntrantsOut(ctx, redEventSlug)
	if err != nil {
		return err
	}
	var redemptionNames []string
	for _, entrant := range redEntrantsResp.Event.Entrants.Nodes {
		redemptionNames = append(redemptionNames, entrant.Name)
	}

	eventSlug := fmt.Sprintf(startgg.EventSlugFormat, tournamentSlug, startgg.EventUltimateSingles)
	entrantsResp, err := startgg.GetEntrantsOut(ctx, eventSlug)
	if err != nil {
		return err
	}
	entrants := entrantsResp.Event.Entrants.Nodes

	if len(entrants) == 0 {
		log.Error("there are no entrants that are out of main bracket and not in redemption")
		return nil
	}

	for _, entrant := range entrants {
		if entrant.Standing.IsFinal && int64(entrant.Standing.Placement) > cutoff && !slices.Contains(redemptionNames, entrant.Name) {
			player := entrant.Participants[0].Player
			startgg.PrintPlayerSimple(player.GamerTag, player.Id)
		}
	}
	return nil
}
