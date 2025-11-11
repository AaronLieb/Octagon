package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestRandomizeSeeds(t *testing.T) {
	players := []brackets.Player{
		{ID: 1, Name: "Player1"},
		{ID: 2, Name: "Player2"},
		{ID: 3, Name: "Player3"},
		{ID: 4, Name: "Player4"},
		{ID: 5, Name: "Player5"},
	}

	// Test that top 2 seeds are preserved
	for range 10 {
		randomized := randomizeSeeds(players, 2)

		if randomized[0].ID != 1 {
			t.Error("Expected first seed to remain unchanged")
		}
		if randomized[1].ID != 2 {
			t.Error("Expected second seed to remain unchanged")
		}

		// Should have same length
		if len(randomized) != len(players) {
			t.Error("Expected same number of players after randomization")
		}
	}
}

func TestResolveConflictsNoConflicts(t *testing.T) {
	players := []brackets.Player{
		{ID: 1, Name: "Player1"},
		{ID: 2, Name: "Player2"},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
		},
	}

	// No conflicts - should return original players
	result := ResolveConflicts(bracket, []Conflict{}, players)

	if len(result) != len(players) {
		t.Error("Expected same number of players")
	}

	// Should be identical (no optimization needed)
	for i, player := range result {
		if player.ID != players[i].ID {
			t.Error("Expected no changes when no conflicts exist")
		}
	}
}
