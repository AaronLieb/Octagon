package conflicts

import (
	"math"

	"github.com/AaronLieb/octagon/brackets"
)

const (
	pruningThreshold = 20.0
	pruningPenalty   = 1000.0
)

// calculateConflictScoreCached is optimized version using pre-computed cache
func calculateConflictScoreCached(cache *conflictCache, conflicts []Conflict, newPlayers []brackets.Player) float64 {
	conflictScore, _ := checkConflictCached(cache, conflicts, newPlayers)

	// Early pruning
	if conflictScore > pruningThreshold {
		return conflictScore + pruningPenalty
	}

	return conflictScore + calculateSeedDiffScoreCached(cache, newPlayers)
}

// calculateSeedDiffScore computes penalty for seed position changes
func calculateSeedDiffScore(players, newPlayers []brackets.Player) float64 {
	seedDiffScore := 0.0
	for i, p := range newPlayers {
		for j, q := range players {
			if p.Id == q.Id {
				importance := calculateImportance(j)
				seedDiff := math.Abs(float64(i-j)) * importance
				seedDiffScore += math.Pow(seedDiff, 1.5)
			}
		}
	}
	return seedDiffScore / 2
}

// calculateSeedDiffScoreCached is optimized version using index map
func calculateSeedDiffScoreCached(cache *conflictCache, newPlayers []brackets.Player) float64 {
	seedDiffScore := 0.0
	for i, p := range newPlayers {
		if originalIndex, exists := cache.playerIndexMap[p.Id]; exists {
			importance := calculateImportance(originalIndex)
			seedDiff := math.Abs(float64(i-originalIndex)) * importance
			seedDiffScore += math.Pow(seedDiff, 1.5)
		}
	}
	return seedDiffScore / 2
}

// calculateImportance determines how important it is to maintain a seed position
func calculateImportance(seedIndex int) float64 {
	importance := 32.0 / math.Pow(math.Log2(float64(seedIndex+1)), 3)
	return math.Min(1, math.Max(0.25, importance))
}
