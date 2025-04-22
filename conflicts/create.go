package conflicts

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

func CreateConflictsForSetsPlayed(ctx context.Context, eventSlug string) []conflict {
	log.Debug("create conflicts for set played")

	var conflicts []conflict

	setsResp, err := startgg.GetSets(ctx, eventSlug)
	if err != nil || len(setsResp.Event.Sets.Nodes) == 0 {
		log.Errorf("unable to fetch sets for '%s'", eventSlug)
		return conflicts
	}
	sets := setsResp.Event.Sets.Nodes
	log.Debug("sets", "n", len(sets))

	for _, set := range sets {
		p1 := set.Slots[0].Entrant.Participants[0].Player
		p2 := set.Slots[1].Entrant.Participants[0].Player
		log.Debug("creating recently played conflict", "p1", p1.GamerTag, "p2", p2.GamerTag)
		conflicts = append(conflicts, conflict{
			Priority: 2,
			Players: []player{
				{
					Name: p1.GamerTag,
					Id:   p1.Id,
				},
				{
					Name: p2.GamerTag,
					Id:   p2.Id,
				},
			},
			Reason: fmt.Sprintf("recently played in %s", eventSlug),
		})
	}

	return conflicts
}
