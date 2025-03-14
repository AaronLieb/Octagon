package conflicts

import (
	"context"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

func CreateConflictsForSetsPlayed(ctx context.Context, eventSlug string) []conflict {
	log.Debug("create conflicts for set played")

	var conflicts []conflict

	setsResp, err := startgg.GetSets(ctx, eventSlug)
	if err != nil {
		log.Errorf("unable to fetch sets for '%s'", eventSlug)
		return conflicts
	}
	sets := setsResp.Event.Sets.Nodes
	log.Debug("sets", "n", len(sets))

	conflicts = make([]conflict, len(sets))
	for i, set := range sets {
		p1 := set.Slots[0].Entrant.Participants[0].Player
		p2 := set.Slots[1].Entrant.Participants[0].Player
		log.Debug("creating recently played conflict", "p1", p1.GamerTag, "p2", p2.GamerTag)
		conflicts[i] = conflict{
			Priority: 2,
			Players: []player{
				{
					Name: p1.GamerTag,
					Id:   p1.Id,
				},
				{
					Name: p1.GamerTag,
					Id:   p2.Id,
				},
			},
		}
	}

	return conflicts
}
