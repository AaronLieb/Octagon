package conflicts

import (
	"os"
	"testing"
	"time"
)

func TestRemoveExpiredConflicts(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	conflicts := []Conflict{
		{Reason: "active", Expiration: nil},     // Never expires
		{Reason: "expired", Expiration: &past},  // Expired
		{Reason: "future", Expiration: &future}, // Active
	}

	active := removeExpiredConflicts(conflicts)

	if len(active) != 2 {
		t.Errorf("Expected 2 active conflicts, got %d", len(active))
	}

	// Check that expired conflict was removed
	for _, conflict := range active {
		if conflict.Reason == "expired" {
			t.Error("Expired conflict should have been removed")
		}
	}
}

func TestWriteReadConflictsFile(t *testing.T) {
	conflicts := []Conflict{
		{
			Priority: 1,
			Reason:   "test conflict",
			Players: []Player{
				{Name: "Player1", ID: int64(123)},
				{Name: "Player2", ID: int64(456)},
			},
		},
	}

	// Write to temp file
	tempFile := "test_conflicts.json"

	err := writeConflictsFile(tempFile, conflicts)
	if err != nil {
		t.Fatalf("Failed to write conflicts file: %v", err)
	}

	// Read back
	readConflicts := readConflictsFile(tempFile)

	if len(readConflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(readConflicts))
	}

	if readConflicts[0].Reason != "test conflict" {
		t.Error("Conflict reason not preserved")
	}

	if len(readConflicts[0].Players) != 2 {
		t.Error("Expected 2 players in conflict")
	}

	_ = os.Remove(tempFile)
}
