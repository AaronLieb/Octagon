package seed

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/AaronLieb/octagon/ratings"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "seed",
		Usage:   "Seed a bracket",
		Aliases: []string{"s"},
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

// TODO: add seed change to this and move printing to seed.go
// TODO: add aliases/main account
type player struct {
	name   string
	id     int
	rating float64
}

func seed(ctx context.Context, cmd *cli.Command) error {
	log.Debug("seed")

	tournamentShortSlug := cmd.String("tournament")

	tournamentResp, err := startgg.GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	tournamentSlug := tournamentResp.Tournament.Slug

	event := "ultimate-singles"
	if cmd.Bool("redemption") {
		event = "redemption-bracket"
	}

	slug := fmt.Sprintf("%s/event/%s", tournamentSlug, event)
	log.Debug(slug)
	entrantsResp, err := startgg.GetUsersForEvent(ctx, slug)
	if err != nil {
		return err
	}
	entrants := entrantsResp.Event.Entrants.Nodes

	players := make([]player, len(entrants))
	for i, entrant := range entrants {
		id := entrant.Participants[0].Player.Id
		name := entrant.Participants[0].GamerTag
		log.Debug("Seed", "name", name, "id", id)
		rating, err := ratings.Get(ctx, id)
		if err != nil {
			log.Errorf("error while trying to fetch rating for '%s': %v", name, err)
		}
		if rating == 0.0 {
			log.Warnf("unable to fetch rating for '%s'", name)
		}
		players[i] = player{
			name:   name,
			id:     id,
			rating: rating,
		}
	}

	slices.SortFunc(players, func(a, b player) int {
		return cmp.Compare(b.rating, a.rating)
	})

	bracket := createBracket(len(players))
	conflictFiles := []string{
		CONFLICT_FILE,
	}
	if cmd.String("file") != "" {
		conflictFiles = append(conflictFiles, cmd.String("file"))
	}
	conflicts := getConflicts(conflictFiles)
	resolveConflicts(bracket, conflicts, players)

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

func publishSeeds(ctx context.Context, eventSlug string, players []player) error {
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
			if seed.Players[0].Id == player.id {
				seedMapping[s] = startgg.UpdatePhaseSeedInfo{
					SeedId:  seed.Id,
					SeedNum: s + 1,
				}
			}
		}
	}

	slices.SortFunc(seedMapping, func(a, b startgg.UpdatePhaseSeedInfo) int {
		return cmp.Compare(a.SeedNum, b.SeedNum)
	})

	updateSeedsResp, err := startgg.UpdateSeeding(ctx, phase.Id, seedMapping)
	if err != nil {
		return err
	}
	log.Debug("updateSeedsResp", "phaseId", updateSeedsResp.UpdatePhaseSeeding.Id)

	log.Infof("Successfully seeded %d players for %s", len(players), eventSlug)

	return nil
}
