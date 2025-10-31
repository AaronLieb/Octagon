package bracket

import (
	"context"
	"fmt"
	"sort"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func seedsCommand() *cli.Command {
	return &cli.Command{
		Name:    "seeds",
		Usage:   "List all seeds sorted by highest to lowest",
		Aliases: []string{"ls"},
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
				Usage:   "The event name",
				Value:   startgg.EventUltimateSingles,
			},
		},
		Action: listSeeds,
	}
}

func listSeeds(ctx context.Context, cmd *cli.Command) error {
	log.Debug("list seeds")

	tournamentShortSlug := cmd.String("tournament")
	event := cmd.String("event")

	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	slug := fmt.Sprintf(startgg.EventSlugFormat, tournamentSlug, event)
	log.Debug(slug)

	var allSeeds []startgg.GetSeedsPaginatedEventPhasesPhaseSeedsSeedConnectionNodesSeed
	page := 1

	for {
		seedsResp, err := startgg.GetSeedsPaginated(ctx, slug, page)
		if err != nil {
			return err
		}

		if len(seedsResp.Event.Phases) == 0 {
			break
		}

		seeds := seedsResp.Event.Phases[0].Seeds.Nodes
		if len(seeds) == 0 {
			break
		}

		allSeeds = append(allSeeds, seeds...)
		page++
	}

	sort.Slice(allSeeds, func(i, j int) bool {
		return allSeeds[i].SeedNum < allSeeds[j].SeedNum
	})

	for _, seed := range allSeeds {
		startgg.PrintSeed(seed.Players[0].GamerTag, seed.SeedNum)
	}

	return nil
}
