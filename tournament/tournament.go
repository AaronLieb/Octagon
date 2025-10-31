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

func FetchReportableSets(ctx context.Context, eventSlug string) ([]Set, error) {
	resp, err := startgg.GetReportableSets(ctx, eventSlug)
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
		})
	}

	log.Debugf("Found %d reportable sets", len(sets))
	return sets, nil
}

func parseRound(round int) string {
	if round < 0 {
		return fmt.Sprintf("LR%d", -round)
	}
	return fmt.Sprintf("WR%d", round)
}

func ReportSet(ctx context.Context, set Set, gameResults []GameResult) error {
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
