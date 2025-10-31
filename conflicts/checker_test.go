package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestConflictCheck(t *testing.T) {
	conflict := Conflict{
		Priority: 1,
		Reason:   "test conflict",
		Players: []Player{
			{Name: "Player1", ID: int64(123)},
			{Name: "Player2", ID: int64(456)},
		},
	}

	// Test matching conflict
	if !conflict.check(int64(123), int64(456)) {
		t.Error("Expected conflict to match players 123 and 456")
	}

	// Test non-matching conflict
	if conflict.check(int64(123), int64(789)) {
		t.Error("Expected no conflict between players 123 and 789")
	}

	// Test same player
	if conflict.check(int64(123), int64(123)) {
		t.Error("Expected no conflict for same player")
	}
}

func TestCheckConflict(t *testing.T) {
	players := []brackets.Player{
		{Id: int64(123), Name: "Player1"},
		{Id: int64(456), Name: "Player2"},
		{Id: int64(789), Name: "Player3"},
	}

	conflicts := []Conflict{
		{
			Priority: 1,
			Players: []Player{
				{ID: int64(123)}, {ID: int64(456)},
			},
		},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2}, // Players 123 vs 456 - should conflict
			{Player1: 2, Player2: 3}, // Players 456 vs 789 - no conflict
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	
	if count != 1 {
		t.Errorf("Expected 1 conflict, got %d", count)
	}
	
	if score != 3.0 { // 2 + priority(1)
		t.Errorf("Expected score 3.0, got %.1f", score)
	}
}
