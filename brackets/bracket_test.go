package brackets

import (
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
	result := reduce(round)
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
		numPlayers int
		minSets    int
	}{
		{8, 7},  // 8-player bracket needs at least 7 sets
		{4, 3},  // 4-player bracket needs at least 3 sets
		{2, 1},  // 2-player bracket needs 1 set
	}
	
	for _, test := range tests {
		bracket := CreateBracket(test.numPlayers)
		
		if bracket == nil {
			t.Errorf("CreateBracket(%d) returned nil", test.numPlayers)
			continue
		}
		
		if len(bracket.Sets) < test.minSets {
			t.Errorf("CreateBracket(%d): expected at least %d sets, got %d", 
				test.numPlayers, test.minSets, len(bracket.Sets))
		}
		
		if len(bracket.WinnersRounds) == 0 {
			t.Errorf("CreateBracket(%d): no winners rounds created", test.numPlayers)
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
	
	result := carryDown(lr, wr, false)
	
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
