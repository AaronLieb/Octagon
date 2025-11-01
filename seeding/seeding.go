package seeding

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/conflicts"
	"github.com/AaronLieb/octagon/ratings"
	"github.com/AaronLieb/octagon/startgg"
)

func filterOutNonNumbers(input string) string {
	reg, _ := regexp.Compile("[0-9]+")
	return reg.FindString(input)
}

// GenerateSeeding creates the seeding order with conflict resolution
func GenerateSeeding(ctx context.Context, tournamentSlug string, redemption bool, conflictFiles []string) ([]brackets.Player, error) {
	fullTournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentSlug)
	if err != nil {
		return nil, err
	}

	tournamentSeriesNumber, err := strconv.Atoi(filterOutNonNumbers(fullTournamentSlug))
	if err != nil {
		// Continue without conflict generation if we can't determine series number
	}

	event := startgg.EventUltimateSingles
	consSeriesNumber := tournamentSeriesNumber - 1

	if redemption {
		event = startgg.EventRedemptionBracket
		consSeriesNumber = tournamentSeriesNumber
	}

	consEvent := fmt.Sprintf(startgg.TournamentEventSlugFormat, fmt.Sprintf("octagon-%d", consSeriesNumber), startgg.EventUltimateSingles)
	cons := conflicts.CreateConflictsForSetsPlayed(ctx, consEvent)
	cons = append(cons, conflicts.GetConflicts(conflictFiles)...)

	slug := fmt.Sprintf(startgg.EventSlugFormat, fullTournamentSlug, event)
	entrantsResp, err := startgg.GetUsersForEvent(ctx, slug)
	if err != nil {
		return nil, err
	}
	entrants := entrantsResp.Event.Entrants.Nodes

	players := make([]brackets.Player, len(entrants))
	for i, entrant := range entrants {
		id := entrant.Participants[0].Player.Id
		name := entrant.Participants[0].GamerTag
		rating, err := ratings.Get(ctx, id)
		if err != nil {
			// Continue with 0 rating if fetch fails
		}
		players[i] = brackets.Player{
			Name:   name,
			ID:     id,
			Rating: rating,
		}
	}

	slices.SortFunc(players, func(a, b brackets.Player) int {
		return cmp.Compare(b.Rating, a.Rating)
	})

	bracket := brackets.CreateBracket(len(players))
	players = conflicts.ResolveConflicts(bracket, cons, players)

	return players, nil
}

// PublishSeeding updates the seeding on start.gg
func PublishSeeding(ctx context.Context, eventSlug string, players []brackets.Player) error {
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
			if startgg.ToString(seed.Players[0].Id) == startgg.ToString(player.ID) {
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

	_, err = startgg.UpdateSeeding(ctx, phase.Id, seedMapping)
	return err
}
