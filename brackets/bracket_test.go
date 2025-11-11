package brackets

import (
	"fmt"
	"testing"
)

func TestCreateRound(t *testing.T) {
	tests := []struct {
		n        int
		k        int
		expected []int
	}{
		{2, 0, []int{1, 2}},
		{4, 0, []int{1, 4, 2, 3}},
		{8, 0, []int{1, 8, 4, 5, 2, 7, 3, 6}},
		{1, 0, []int{}},
		{0, 0, []int{}},
	}

	for _, test := range tests {
		result := CreateRound(test.n, test.k)
		if len(result) != len(test.expected) {
			t.Errorf("CreateRound(%d, %d): expected length %d, got %d",
				test.n, test.k, len(test.expected), len(result))
			continue
		}

		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("CreateRound(%d, %d): expected %v, got %v",
					test.n, test.k, test.expected, result)
				break
			}
		}
	}
}

func TestCreateSets(t *testing.T) {
	round := []int{1, 8, 4, 5}
	sets := createSets(round)

	if len(sets) != 2 {
		t.Errorf("Expected 2 sets, got %d", len(sets))
	}

	if sets[0].Player1 != 1 || sets[0].Player2 != 8 {
		t.Errorf("First set: expected players 1 vs 8, got %d vs %d",
			sets[0].Player1, sets[0].Player2)
	}

	if sets[1].Player1 != 4 || sets[1].Player2 != 5 {
		t.Errorf("Second set: expected players 4 vs 5, got %d vs %d",
			sets[1].Player1, sets[1].Player2)
	}
}

func TestReduce(t *testing.T) {
	round := []int{1, 8, 4, 5, 2, 7, 3, 6}
	result := reduceWinners(round)
	expected := []int{1, 4, 2, 3}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
			break
		}
	}
}

func TestCreateBracket(t *testing.T) {
	tests := []struct {
		numPlayers       int
		numSets          int
		numWinnersRounds int
		numLosersRounds  int
	}{
		{16, 29, 4, 6},
		{8, 13, 3, 4},
		{4, 5, 2, 2},
	}

	for _, test := range tests {
		bracket := CreateBracket(test.numPlayers)

		if bracket == nil {
			t.Errorf("CreateBracket(%d) returned nil", test.numPlayers)
			continue
		}

		if len(bracket.Sets) != test.numSets {
			t.Errorf("CreateBracket(%d): expected  %d sets, got %d",
				test.numPlayers, test.numSets, len(bracket.Sets))
		}

		if len(bracket.WinnersRounds) != test.numWinnersRounds {
			t.Errorf("CreateBracket(%d): expected %d winners rounds, got %d", test.numPlayers, test.numWinnersRounds, len(bracket.WinnersRounds))
		}

		if len(bracket.LosersRounds) != test.numLosersRounds {
			t.Errorf("CreateBracket(%d): expected %d losers rounds, got %d", test.numPlayers, test.numLosersRounds, len(bracket.LosersRounds))
		}

		// Verify no invalid player numbers
		for _, set := range bracket.Sets {
			if set.Player1 > test.numPlayers || set.Player2 > test.numPlayers {
				t.Errorf("CreateBracket(%d): invalid player numbers %d vs %d",
					test.numPlayers, set.Player1, set.Player2)
			}
		}
	}
}

func TestCarryDown(t *testing.T) {
	lr := []int{5, 6}
	wr := []int{1, 2}

	result := carryDown(lr, wr, 8, 0)

	if len(result) != len(lr) {
		t.Errorf("Expected length %d, got %d", len(lr), len(result))
	}

	// Basic structure test - should have proper length
	if len(result) == 0 {
		t.Error("carryDown returned empty result")
	}
}

func TestBracketStructure(t *testing.T) {
	bracket := CreateBracket(8)

	// Test that bracket has proper structure
	if bracket.Sets == nil {
		t.Error("Bracket.Sets is nil")
	}

	if bracket.WinnersRounds == nil {
		t.Error("Bracket.WinnersRounds is nil")
	}

	if bracket.LosersRounds == nil {
		t.Error("Bracket.LosersRounds is nil")
	}

	// Test that all sets have valid player assignments
	for i, set := range bracket.Sets {
		if set.Player1 <= 0 || set.Player2 <= 0 {
			t.Errorf("Set %d has invalid player numbers: %d vs %d",
				i, set.Player1, set.Player2)
		}
	}
}

func TestBracket128PlayerSets(t *testing.T) {
	bracket := CreateBracket(128)

	expectedSets := []struct {
		round   int
		player1 int
		player2 int
	}{
		{1, 25, 104},
		{1, 61, 68},
		{-1, 82, 111},
		{-1, 87, 106},
		{-1, 78, 115},
		{-2, 61, 87},
		{-2, 48, 70},
		{-3, 42, 55},
		{-3, 43, 54},
		{-3, 33, 64},
		{-4, 18, 48},
		{-5, 23, 26},
		{-5, 22, 27},
		{-6, 11, 17},
		{-7, 11, 14},
		{-8, 8, 10},
		{-9, 6, 7},
		{-10, 3, 5},
		{-11, 3, 4},
	}

	for _, expected := range expectedSets {
		found := false
		var sets []*Set
		if expected.round > 0 {
			sets = bracket.WinnersRounds[expected.round-1]
		} else {
			sets = bracket.LosersRounds[-expected.round-1]
		}
		for _, set := range sets {
			if (set.Player1 == expected.player1 && set.Player2 == expected.player2) ||
				(set.Player1 == expected.player2 && set.Player2 == expected.player1) {
				found = true
				break
			}
		}
		if !found {
			setsFound := ""
			for _, set := range sets {
				setsFound += fmt.Sprintf("[%d,%d]", set.Player1, set.Player2)
			}
			t.Errorf("128-player bracket missing expected set %d: %d vs %d, found: {%s}",
				expected.round, expected.player1, expected.player2, setsFound)
		}
	}
}

func TestBracket64PlayerSets(t *testing.T) {
	bracket := CreateBracket(64)

	expectedSets := []struct {
		round   int
		player1 int
		player2 int
	}{
		{-1, 48, 49},
		{-2, 30, 41},
		{-3, 23, 26},
		{-4, 11, 21},
		{-5, 10, 15},
		{-6, 8, 11},
		{-7, 6, 7},
		{-7, 8, 5},
		{-8, 3, 5},
	}

	for _, expected := range expectedSets {
		found := false
		var sets []*Set
		if expected.round > 0 {
			sets = bracket.WinnersRounds[expected.round-1]
		} else {
			sets = bracket.LosersRounds[-expected.round-1]
		}
		for _, set := range sets {
			if (set.Player1 == expected.player1 && set.Player2 == expected.player2) ||
				(set.Player1 == expected.player2 && set.Player2 == expected.player1) {
				found = true
				break
			}
		}
		if !found {
			setsFound := ""
			for _, set := range sets {
				setsFound += fmt.Sprintf("[%d,%d]", set.Player1, set.Player2)
			}
			t.Errorf("64-player bracket missing expected set %d: %d vs %d, found: {%s}",
				expected.round, expected.player1, expected.player2, setsFound)
		}
	}
}

func TestBracket32PlayerSets(t *testing.T) {
	bracket := CreateBracket(32)

	expectedSets := []struct {
		round   int
		player1 int
		player2 int
	}{
		{1, 16, 17},
		{2, 4, 13},
		{3, 2, 7},
		{4, 1, 4},
		{5, 1, 2},
		{-1, 28, 21},
		{-1, 26, 23},
		{-2, 11, 17},
		{-2, 13, 23},
		{-2, 9, 19},
		{-3, 10, 15},
		{-3, 9, 16},
		{-4, 5, 11},
		{-4, 6, 12},
		{-5, 5, 8},
		{-6, 3, 5},
		{-7, 3, 4},
		{-8, 2, 3},
	}

	for _, expected := range expectedSets {
		found := false
		var sets []*Set
		if expected.round > 0 {
			sets = bracket.WinnersRounds[expected.round-1]
		} else {
			sets = bracket.LosersRounds[-expected.round-1]
		}
		for _, set := range sets {
			if (set.Player1 == expected.player1 && set.Player2 == expected.player2) ||
				(set.Player1 == expected.player2 && set.Player2 == expected.player1) {
				found = true
				break
			}
		}
		if !found {
			setsFound := ""
			for _, set := range sets {
				setsFound += fmt.Sprintf("[%d,%d]", set.Player1, set.Player2)
			}
			t.Errorf("32-player bracket missing expected set %d: %d vs %d. Found: {%s}",
				expected.round, expected.player1, expected.player2, setsFound)
		}
	}
}
