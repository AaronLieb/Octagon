package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestNewConflictCache(t *testing.T) {
	players := []brackets.Player{
		{Id: int64(123), Name: "Player1"},
		{Id: int64(456), Name: "Player2"},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
			nil, // Should be filtered out
			{Player1: 2, Player2: 1},
		},
	}

	cache := newConflictCache(players, bracket)

	// Test player index map
	if cache.playerIndexMap[int64(123)] != 0 {
		t.Error("Expected player 123 at index 0")
	}
	if cache.playerIndexMap[int64(456)] != 1 {
		t.Error("Expected player 456 at index 1")
	}

	// Test string cache
	if cache.stringCache[int64(123)] != "123" {
		t.Error("Expected string cache for player 123")
	}

	// Test conflict sets (should filter out nil)
	if len(cache.conflictSets) != 2 {
		t.Errorf("Expected 2 conflict sets, got %d", len(cache.conflictSets))
	}
}

func TestCheckCached(t *testing.T) {
	players := []brackets.Player{
		{Id: int64(123), Name: "Player1"},
		{Id: int64(456), Name: "Player2"},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
		},
	}

	cache := newConflictCache(players, bracket)

	conflict := Conflict{
		Players: []Player{
			{ID: int64(123)}, {ID: int64(456)},
		},
	}

	// Test matching
	if !conflict.checkCached(int64(123), int64(456), cache) {
		t.Error("Expected cached conflict check to match")
	}

	// Test fallback for uncached ID
	conflict.Players[0].ID = int64(999) // Not in cache
	// This should work because checkCached falls back to ToString for uncached IDs
	if conflict.checkCached(int64(999), int64(456), cache) {
		t.Error("Expected no match when one player is not in bracket")
	}
}
