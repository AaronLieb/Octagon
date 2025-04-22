package startgg

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

type Set struct {
	Player1 brackets.Player
	Player2 brackets.Player
	Round   string
}

func FetchReportableSets(ctx context.Context, eventSlug string) []*Set {
	setsResp, err := startgg.GetReportableSets(ctx, eventSlug)
	if err != nil {
		log.Fatalf("unable to find sets: %v", err)
	}
	setNodes := setsResp.Event.Sets.Nodes
	var sets []*Set
	for _, set := range setNodes {
		parts1 := set.Slots[0].Entrant.Participants
		parts2 := set.Slots[1].Entrant.Participants
		if len(parts1) > 0 && len(parts2) > 0 {
			p1 := set.Slots[0].Entrant.Participants[0].Player
			p2 := set.Slots[1].Entrant.Participants[0].Player
			sets = append(sets, &Set{
				Player1: brackets.Player{
					Name: p1.GamerTag,
					Id:   p1.Id,
				},
				Player2: brackets.Player{
					Name: p2.GamerTag,
					Id:   p2.Id,
				},
				Round: parseRound(set.Round),
			})
		}
	}
	return sets
}

func parseRound(round int) string {
	if round < 0 {
		return fmt.Sprintf("LR%d", -round)
	} else {
		return fmt.Sprintf("WR%d", round)
	}
}
