package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestRandomizeSeeds(t *testing.T) {
	players := []brackets.Player{
		{Id: 1, Name: "Player1"},
		{Id: 2, Name: "Player2"},
		{Id: 3, Name: "Player3"},
		{Id: 4, Name: "Player4"},
		{Id: 5, Name: "Player5"},
	}

	// Test that top 2 seeds are preserved
	for i := 0; i < 10; i++ {
		randomized := randomizeSeeds(players, 2)
		
		if randomized[0].Id != 1 {
			t.Error("Expected first seed to remain unchanged")
		}
		if randomized[1].Id != 2 {
			t.Error("Expected second seed to remain unchanged")
		}
		
		// Should have same length
		if len(randomized) != len(players) {
			t.Error("Expected same number of players after randomization")
		}
	}
}

func TestCalculateAttempts(t *testing.T) {
	tests := []struct {
		variance int
		expected int
	}{
		{6, ConflictResolutionAttempts},
		{5, int(ConflictResolutionAttempts * 1.0)},
		{3, int(ConflictResolutionAttempts * 5.196)}, // 6-3 = 3, 3^1.5 ≈ 5.196
	}

	for _, test := range tests {
		result := calculateAttempts(test.variance)
		if result < ConflictResolutionAttempts {
			t.Errorf("calculateAttempts(%d) = %d, should be at least %d", 
				test.variance, result, ConflictResolutionAttempts)
		}
	}
}

func TestResolveConflictsNoConflicts(t *testing.T) {
	players := []brackets.Player{
		{Id: 1, Name: "Player1"},
		{Id: 2, Name: "Player2"},
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
		if player.Id != players[i].Id {
			t.Error("Expected no changes when no conflicts exist")
		}
	}
}
