package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestCalculateImportance(t *testing.T) {
	// Test that importance is between 0.25 and 1.0
	for i := 0; i < 32; i++ {
		result := calculateImportance(i)
		if result < 0.25 || result > 1.0 {
			t.Errorf("calculateImportance(%d) = %.2f, should be between 0.25 and 1.0", i, result)
		}
	}
	
	// Test that top seeds have higher importance than lower seeds
	top := calculateImportance(0)
	mid := calculateImportance(15)
	low := calculateImportance(31)
	
	if top <= mid || mid <= low {
		t.Error("Expected decreasing importance for lower seeds")
	}
}

func TestCalculateSeedDiffScore(t *testing.T) {
	original := []brackets.Player{
		{Id: int64(1), Name: "Player1"},
		{Id: int64(2), Name: "Player2"},
		{Id: int64(3), Name: "Player3"},
	}

	// No changes
	same := []brackets.Player{
		{Id: int64(1), Name: "Player1"},
		{Id: int64(2), Name: "Player2"},
		{Id: int64(3), Name: "Player3"},
	}

	score := calculateSeedDiffScore(original, same)
	if score != 0.0 {
		t.Errorf("Expected 0 score for identical seeding, got %.2f", score)
	}

	// Swap positions
	swapped := []brackets.Player{
		{Id: int64(2), Name: "Player2"},
		{Id: int64(1), Name: "Player1"},
		{Id: int64(3), Name: "Player3"},
	}

	score = calculateSeedDiffScore(original, swapped)
	if score <= 0.0 {
		t.Error("Expected positive score for swapped seeding")
	}
}
