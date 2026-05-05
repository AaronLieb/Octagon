package conflicts

import (
	"math"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/startgg"
)

// check returns true if p1 and p2 are in the conflict
func (con *Conflict) check(p1, p2 startgg.ID) bool {
	flag := true
	for _, p := range con.Players {
		if startgg.ToString(p.ID) == startgg.ToString(p1) || startgg.ToString(p.ID) == startgg.ToString(p2) {
			if flag {
				flag = false
			} else {
				return true
			}
		}
	}
	return false
}

// checkCached is optimized version using cached string lookups
func (con *Conflict) checkCached(p1, p2 startgg.ID, cache *conflictCache) bool {
	p1Str := cache.stringCache[p1]
	p2Str := cache.stringCache[p2]

	flag := true
	for _, p := range con.Players {
		pStr := cache.stringCache[p.ID]
		// Fallback to ToString if not in cache
		if pStr == "" {
			pStr = startgg.ToString(p.ID)
		}

		if pStr == p1Str || pStr == p2Str {
			if flag {
				flag = false
			} else {
				return true
			}
		}
	}
	return false
}

// isRequest returns true if the conflict is a set request (negative priority)
func (con *Conflict) isRequest() bool {
	return con.Priority < 0
}

// setMatchesRound checks if a set is in the conflict's target round.
// Returns true if the conflict has no round restriction.
func setMatchesRound(s *brackets.Set, roundMap map[*brackets.Set]brackets.Round, con *Conflict) bool {
	if con.Round == nil {
		return true
	}
	if r, ok := roundMap[s]; ok {
		return r == *con.Round
	}
	return false
}

// setsForRound returns the subset of sets that match the conflict's round.
// If Round is nil, returns all sets.
func setsForRound(sets []*brackets.Set, roundMap map[*brackets.Set]brackets.Round, con *Conflict) []*brackets.Set {
	if con.Round == nil {
		return sets
	}
	var filtered []*brackets.Set
	for _, s := range sets {
		if r, ok := roundMap[s]; ok && r == *con.Round {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// buildRoundMap creates a set-to-round mapping from a bracket.
// Round numbers are 1-indexed to match brackets.ParseRound, Round.String,
// and brackets.RoundFromStartGG — so the first winners round is WR1, not WR0.
func buildRoundMap(bracket *brackets.Bracket) map[*brackets.Set]brackets.Round {
	roundMap := make(map[*brackets.Set]brackets.Round)
	for i, sets := range bracket.WinnersRounds {
		for _, s := range sets {
			roundMap[s] = brackets.Round{Number: i + 1}
		}
	}
	for i, sets := range bracket.LosersRounds {
		for _, s := range sets {
			roundMap[s] = brackets.Round{Losers: true, Number: i + 1}
		}
	}
	return roundMap
}

// checkConflict returns the conflict sum and count for a bracket
func checkConflict(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) (float64, int) {
	roundMap := buildRoundMap(bracket)

	conflictScore := 0.0
	conflictSum := 0

	// Score regular conflicts: penalize when pair IS found
	for _, s := range bracket.Sets {
		if s == nil {
			continue
		}
		for _, con := range conflicts {
			if con.isRequest() {
				continue
			}
			if !setMatchesRound(s, roundMap, &con) {
				continue
			}
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.check(p1, p2) {
				conflictScore += float64(2 + con.Priority)
				conflictSum++
			}
		}
	}

	// Score set requests: penalize when pair is NOT found
	for i := range conflicts {
		con := &conflicts[i]
		if !con.isRequest() {
			continue
		}
		found := false
		for _, s := range setsForRound(bracket.Sets, roundMap, con) {
			if s == nil {
				continue
			}
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.check(p1, p2) {
				found = true
				break
			}
		}
		if !found {
			conflictScore += 2 + math.Abs(float64(con.Priority))
			conflictSum++
		}
	}

	return conflictScore, conflictSum
}

// checkConflictCached is an optimized version of checkConflict that reuses
// pre-computed caches. It must return the same (score, count) as checkConflict
// for the same inputs — the cache only speeds up ID-string lookups and
// set-to-round lookups.
func checkConflictCached(cache *conflictCache, conflicts []Conflict, players []brackets.Player) (float64, int) {
	conflictScore := 0.0
	conflictSum := 0

	// Score regular conflicts: penalize when the pair IS found in a matching round.
	for _, s := range cache.conflictSets {
		for i := range conflicts {
			con := &conflicts[i]
			if con.isRequest() {
				continue
			}
			if !setMatchesRound(s, cache.setRoundMap, con) {
				continue
			}
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.checkCached(p1, p2, cache) {
				conflictScore += float64(2 + con.Priority)
				conflictSum++
			}
		}
	}

	// Score set requests: penalize (once per request) when the pair is NOT found
	// in any set matching the request's round (if any).
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
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.checkCached(p1, p2, cache) {
				found = true
				break
			}
		}
		if !found {
			conflictScore += 2 + math.Abs(float64(con.Priority))
			conflictSum++
		}
	}

	return conflictScore, conflictSum
}

// listUnresolvedConflicts returns conflicts and unfulfilled requests
func listUnresolvedConflicts(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) []Conflict {
	roundMap := buildRoundMap(bracket)
	var unresolved []Conflict

	for i := range conflicts {
		con := &conflicts[i]
		if con.isRequest() {
			// Request is unresolved if pair is NOT found in matching sets
			found := false
			for _, s := range setsForRound(bracket.Sets, roundMap, con) {
				if s == nil {
					continue
				}
				p1 := players[s.Player1-1].ID
				p2 := players[s.Player2-1].ID
				if con.check(p1, p2) {
					found = true
					break
				}
			}
			if !found {
				unresolved = append(unresolved, *con)
			}
		} else {
			// Conflict is unresolved if pair IS found in matching sets
			for _, s := range bracket.Sets {
				if s == nil {
					continue
				}
				if !setMatchesRound(s, roundMap, con) {
					continue
				}
				p1 := players[s.Player1-1].ID
				p2 := players[s.Player2-1].ID
				if con.check(p1, p2) {
					unresolved = append(unresolved, *con)
					break
				}
			}
		}
	}

	return unresolved
}
