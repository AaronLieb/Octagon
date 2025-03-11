package attendees

import (
	"context"
	"fmt"
	"slices"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

// out before top8
const CUTOFF = 8

func RedemptionInfoCommand() *cli.Command {
	return &cli.Command{
		Name:    "redemption",
		Usage:   "Redemption info",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
		},
		Action: RedemptionInfo,
	}
}

func RedemptionInfo(ctx context.Context, cmd *cli.Command) error {
	log.Debug("Redemption Info")

	tournamentShortSlug := cmd.String("tournament")

	tournamentResp, err := startgg.GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	tournamentSlug := tournamentResp.Tournament.Slug

	redEventSlug := fmt.Sprintf("%s/event/redemption-bracket", tournamentSlug)
	redEntrantsResp, err := startgg.GetEntrantsOut(ctx, redEventSlug)
	if err != nil {
		return err
	}
	var redemptionNames []string
	for _, entrant := range redEntrantsResp.Event.Entrants.Nodes {
		redemptionNames = append(redemptionNames, entrant.Name)
	}

	eventSlug := fmt.Sprintf("tournament/%s/event/ultimate-singles", tournamentSlug)
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
		if entrant.Standing.IsFinal && entrant.Standing.Placement > CUTOFF && !slices.Contains(redemptionNames, entrant.Name) {
			fmt.Println(entrant.Name)
		}
	}
	return nil
}
