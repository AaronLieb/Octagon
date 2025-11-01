package conflicts

import (
	"encoding/json"
	"os"
	"time"

	"github.com/AaronLieb/octagon/config"
	"github.com/charmbracelet/log"
)

const ConflictFile = "conflicts.json"

// GetConflicts reads conflicts from files and removes expired ones
func GetConflicts(conflictFiles []string) []Conflict {
	conflicts := readConflictsFile(ConflictFile)
	for _, file := range conflictFiles {
		conflicts = append(conflicts, readConflictsFile(file)...)
	}

	// Remove expired conflicts and update main file if needed
	originalCount := len(conflicts)
	conflicts = removeExpiredConflicts(conflicts)
	if len(conflicts) < originalCount {
		err := writeConflictsFile(ConflictFile, conflicts)
		if err != nil {
			log.Error("Unable to remove expired conflicts", "err", err)
		}
	}

	return conflicts
}

// SaveConflict adds a new conflict to the file
func SaveConflict(newConflict Conflict) error {
	existingConflicts := readConflictsFile(ConflictFile)
	existingConflicts = removeExpiredConflicts(existingConflicts)
	existingConflicts = append(existingConflicts, newConflict)
	return writeConflictsFile(ConflictFile, existingConflicts)
}

// WriteConflictsFile overwrites the conflicts file with the provided conflicts
func WriteConflictsFile(conflicts []Conflict) error {
	return writeConflictsFile(ConflictFile, conflicts)
}

// readConflictsFile reads conflicts from a JSON file
func readConflictsFile(fileName string) []Conflict {
	var conflicts []Conflict
	path := config.GetConfigPath() + fileName

	file, err := os.ReadFile(path)
	if err != nil {
		log.Warn("unable to read conflicts file: %v", err)
		return conflicts
	}

	if err := json.Unmarshal(file, &conflicts); err != nil {
		log.Errorf("unable to unmarshal conflicts file: %v", err)
		return conflicts
	}

	log.Info("Reading conflicts", "path", path, "n", len(conflicts))
	if len(conflicts) == 0 {
		log.Warn("No conflicts found in conflict file")
	}

	return conflicts
}

// writeConflictsFile writes conflicts to a JSON file
func writeConflictsFile(fileName string, conflicts []Conflict) error {
	path := config.GetConfigPath() + fileName

	data, err := json.MarshalIndent(conflicts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// removeExpiredConflicts filters out expired conflicts
func removeExpiredConflicts(conflicts []Conflict) []Conflict {
	now := time.Now()
	var active []Conflict

	for _, conflict := range conflicts {
		if conflict.Expiration == nil || conflict.Expiration.After(now) {
			active = append(active, conflict)
		}
	}

	return active
}
