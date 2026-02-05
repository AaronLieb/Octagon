// Package tournament provides tournament set management, fetching, and reporting functionality.
package tournament

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

type Set struct {
	ID       int
	Player1  Player
	Player2  Player
	Round    string
	Entrant1 int
	Entrant2 int
	State    int // 1 = not started, 2 = in progress, 3 = completed
}

type Player struct {
	Name string
	ID   int
}

type GameResult struct {
	Winner   int // 1 or 2
	P1Char   string
	P2Char   string
	P1CharID int
	P2CharID int
}

func FetchReportableSets(ctx context.Context, eventSlug string, includeCompleted bool) ([]Set, error) {
	// State 1 = not started, State 2 = in progress, State 3 = completed
	states := []int{1, 2}
	if includeCompleted {
		states = append(states, 3)
	}

	// Note: The current GraphQL query only fetches 500 sets from page 0
	resp, err := startgg.GetReportableSets(ctx, eventSlug, states)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sets: %w", err)
	}

	var sets []Set
	for _, setNode := range resp.Event.Sets.Nodes {
		if len(setNode.Slots) != 2 {
			continue
		}

		slot1 := setNode.Slots[0]
		slot2 := setNode.Slots[1]

		if slot1.Entrant.Id == nil || slot2.Entrant.Id == nil ||
			len(slot1.Entrant.Participants) == 0 || len(slot2.Entrant.Participants) == 0 {
			continue
		}

		p1 := slot1.Entrant.Participants[0].Player
		p2 := slot2.Entrant.Participants[0].Player

		// Convert interface{} to int
		setID, ok := setNode.Id.(float64)
		if !ok {
			continue
		}
		entrant1Id, ok := slot1.Entrant.Id.(float64)
		if !ok {
			continue
		}
		entrant2Id, ok := slot2.Entrant.Id.(float64)
		if !ok {
			continue
		}
		player1Id, ok := p1.Id.(float64)
		if !ok {
			continue
		}
		player2Id, ok := p2.Id.(float64)
		if !ok {
			continue
		}

		sets = append(sets, Set{
			ID: int(setID),
			Player1: Player{
				Name: p1.GamerTag,
				ID:   int(player1Id),
			},
			Player2: Player{
				Name: p2.GamerTag,
				ID:   int(player2Id),
			},
			Round:    parseRound(setNode.Round),
			Entrant1: int(entrant1Id),
			Entrant2: int(entrant2Id),
			State:    setNode.State,
		})
	}

	log.Debugf("Found %d reportable sets (limited to first 500 due to pagination)", len(sets))
	if len(resp.Event.Sets.Nodes) == 500 {
		log.Warnf("Fetched exactly 500 sets - there may be more sets available that aren't shown")
	}
	return sets, nil
}

func parseRound(round int) string {
	if round < 0 {
		return fmt.Sprintf("LR%d", -round)
	}
	return fmt.Sprintf("WR%d", round)
}

func ValidateSetScore(gameResults []GameResult) error {
	p1Wins := 0
	p2Wins := 0
	for _, result := range gameResults {
		switch result.Winner {
		case 1:
			p1Wins++
		case 2:
			p2Wins++
		}
	}

	// Validate best-of-3 (2-0, 2-1) or best-of-5 (3-0, 3-1, 3-2)
	if (p1Wins != 2 || p2Wins > 1) && (p2Wins != 2 || p1Wins > 1) &&
		(p1Wins != 3 || p2Wins > 2) && (p2Wins != 3 || p1Wins > 2) {
		return fmt.Errorf("invalid score: must be best-of-3 (2-0, 2-1) or best-of-5 (3-0, 3-1, 3-2)")
	}

	return nil
}

func ReportSet(ctx context.Context, set Set, gameResults []GameResult) error {
	// If set is already completed (state 3), reset it first
	if set.State == 3 {
		log.Infof("Set %d is already completed, resetting before reporting", set.ID)
		_, err := startgg.ResetSet(ctx, set.ID)
		if err != nil {
			log.Errorf("Failed to reset set: %v", err)
			return fmt.Errorf("failed to reset set: %w", err)
		}
	}

	p1Wins := 0
	p2Wins := 0
	for _, result := range gameResults {
		switch result.Winner {
		case 1:
			p1Wins++
		case 2:
			p2Wins++
		}
	}

	var winnerID int
	if p1Wins > p2Wins {
		winnerID = set.Entrant1
	} else {
		winnerID = set.Entrant2
	}

	var gameData []startgg.BracketSetGameDataInput
	for i, result := range gameResults {
		if result.Winner != 0 {
			var selections []startgg.BracketSetGameSelectionInput

			if result.P1CharID != 0 {
				selections = append(selections, startgg.BracketSetGameSelectionInput{
					EntrantId:   set.Entrant1,
					CharacterId: result.P1CharID,
				})
			}
			if result.P2CharID != 0 {
				selections = append(selections, startgg.BracketSetGameSelectionInput{
					EntrantId:   set.Entrant2,
					CharacterId: result.P2CharID,
				})
			}

			gameData = append(gameData, startgg.BracketSetGameDataInput{
				GameNum: i + 1,
				WinnerId: func() int {
					if result.Winner == 1 {
						return set.Entrant1
					}
					return set.Entrant2
				}(),
				Selections: selections,
			})
		}
	}

	_, err := startgg.ReportSet(ctx, set.ID, winnerID, gameData)
	if err != nil {
		log.Errorf("Failed to report set: %v", err)
	}
	return err
}
