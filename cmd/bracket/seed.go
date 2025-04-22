package bracket

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/conflicts"
	"github.com/AaronLieb/octagon/ratings"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const desc = `Will fetch player ratings from an external rating database
and seed the players accordingly. It will then attempt to read all
player conflicts and generate a variation of the original seeding
that minimizes seeding changes while maximizing conflict resolution`

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

	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	// TODO: Don't hard code this, dynamically generate octagon numbers
	event := "ultimate-singles"
	consEvent := "tournament/octagon-108/event/ultimate-singles"
	if cmd.Bool("redemption") {
		event = "redemption-bracket"
		consEvent = "tournament/octagon/event/ultimate-singles"
	}

	cons := conflicts.CreateConflictsForSetsPlayed(ctx, consEvent)

	slug := fmt.Sprintf("%s/event/%s", tournamentSlug, event)
	log.Debug(slug)
	entrantsResp, err := startgg.GetUsersForEvent(ctx, slug)
	if err != nil {
		return err
	}
	entrants := entrantsResp.Event.Entrants.Nodes

	players := make([]brackets.Player, len(entrants))
	for i, entrant := range entrants {
		id := int64(entrant.Participants[0].Player.Id.(float64))
		name := entrant.Participants[0].GamerTag
		log.Debug("Seed", "name", name, "id", id)
		rating, err := ratings.Get(ctx, id)
		if err != nil {
			log.Errorf("error while trying to fetch rating for '%s': %v", name, err)
		}
		if rating == 0.0 {
			log.Warnf("unable to fetch rating for '%s'", name)
		}
		players[i] = brackets.Player{
			Name:   name,
			Id:     id,
			Rating: rating,
		}
	}

	slices.SortFunc(players, func(a, b brackets.Player) int {
		return cmp.Compare(b.Rating, a.Rating)
	})

	bracket := brackets.CreateBracket(len(players))
	var conflictFiles []string
	if cmd.String("file") != "" {
		conflictFiles = append(conflictFiles, cmd.String("file"))
	}
	cons = append(cons, conflicts.GetConflicts(conflictFiles)...)
	conflicts.ResolveConflicts(bracket, cons, players)

	var input string
	fmt.Println("Publish seeding? (y/N)")
	fmt.Scanln(&input)

	if input != "y" && input != "Y" {
		log.Info("Cancelling seeding...")
		return nil
	}

	err = publishSeeds(ctx, slug, players)
	if err != nil {
		return err
	}

	return nil
}

func publishSeeds(ctx context.Context, eventSlug string, players []brackets.Player) error {
	seedsResp, err := startgg.GetSeeds(ctx, eventSlug)
	if err != nil {
		return err
	}
	phases := seedsResp.Event.Phases

	if len(phases) == 0 {
		return errors.New("no phases in event")
	}

	phase := phases[0]
	seeds := phase.Seeds.Nodes

	seedMapping := make([]startgg.UpdatePhaseSeedInfo, len(players))
	for _, seed := range seeds {
		for s, player := range players {
			if startgg.ToString(seed.Players[0].Id) == startgg.ToString(player.Id) {
				seedMapping[s] = startgg.UpdatePhaseSeedInfo{
					SeedId:  seed.Id,
					SeedNum: s + 1,
				}
			}
		}
	}

	slices.SortFunc(seedMapping, func(a, b startgg.UpdatePhaseSeedInfo) int {
		return cmp.Compare(a.SeedNum.(int), b.SeedNum.(int))
	})

	updateSeedsResp, err := startgg.UpdateSeeding(ctx, phase.Id, seedMapping)
	if err != nil {
		return err
	}
	log.Debug("updateSeedsResp", "phaseId", updateSeedsResp.UpdatePhaseSeeding.Id)

	log.Infof("Successfully seeded %d players for %s", len(players), eventSlug)

	return nil
}
