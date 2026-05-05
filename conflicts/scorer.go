package conflicts

import (
	"math"

	"github.com/AaronLieb/octagon/brackets"
)

const (
	// pruningThreshold short-circuits seedDiff computation for candidates whose
	// conflict/request cost alone is already far above the current best. Set
	// high enough that a handful of unfulfilled priority-3 requests (each
	// costing mcRequestPenalty for priority 3) don't trip it immediately.
	pruningThreshold = 500.0
	pruningPenalty   = 10000.0

	// mcRequestWeight scales the set-request penalty in the Monte Carlo cost
	// function so that requests dominate seedDiff. A set request is an explicit
	// TO directive to make two players play — the optimizer should be willing
	// to accept significant seed rearrangement to fulfill it. Regular conflicts
	// keep their natural weight because they're "avoid if possible", not
	// "enforce at all costs".
	mcRequestWeight = 20.0
)

// calculateConflictScoreCached is the Monte Carlo objective function. Lower is
// better. It is intentionally different from checkConflictCached: set-request
// penalties are amplified here so the optimizer treats them as overrides of
// seed quality, which matches the TO's intent when creating a request.
//
// The returned value is not meaningful outside of MC ranking; user-facing
// reporting should use checkConflict/checkConflictCached instead.
func calculateConflictScoreCached(cache *conflictCache, conflicts []Conflict, newPlayers []brackets.Player) float64 {
	cost := 0.0

	// Regular conflicts: pair found in a matching round → natural penalty.
	for _, s := range cache.conflictSets {
		for i := range conflicts {
			con := &conflicts[i]
			if con.isRequest() {
				continue
			}
			if !setMatchesRound(s, cache.setRoundMap, con) {
				continue
			}
			p1 := newPlayers[s.Player1-1].ID
			p2 := newPlayers[s.Player2-1].ID
			if con.checkCached(p1, p2, cache) {
				cost += float64(2 + con.Priority)
			}
		}
	}

	// Set requests: pair NOT found in any matching set → weighted penalty.
	for i := range conflicts {
		con := &conflicts[i]
		if !con.isRequest() {
			continue
		}
		found := false
		for _, s := range cache.conflictSets {
			if !setMatchesRound(s, cache.setRoundMap, con) {
				continue
			}
			p1 := newPlayers[s.Player1-1].ID
			p2 := newPlayers[s.Player2-1].ID
			if con.checkCached(p1, p2, cache) {
				found = true
				break
			}
		}
		if !found {
			cost += (2 + math.Abs(float64(con.Priority))) * mcRequestWeight
		}
	}

	if cost > pruningThreshold {
		return cost + pruningPenalty
	}

	return cost + calculateSeedDiffScoreCached(cache, newPlayers)
}

// calculateSeedDiffScore computes penalty for seed position changes
func calculateSeedDiffScore(players, newPlayers []brackets.Player) float64 {
	seedDiffScore := 0.0
	for i, p := range newPlayers {
		for j, q := range players {
			if p.ID == q.ID {
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
		if originalIndex, exists := cache.playerIndexMap[p.ID]; exists {
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
