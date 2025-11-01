package conflicts

import (
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

// checkConflict returns the conflict sum and count for a bracket
func checkConflict(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) (float64, int) {
	conflictScore := 0.0
	conflictSum := 0

	for _, s := range bracket.Sets {
		if s == nil {
			continue
		}
		for _, con := range conflicts {
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.check(p1, p2) {
				conflictScore += float64(2 + con.Priority)
				conflictSum++
			}
		}
	}

	return conflictScore, conflictSum
}

// checkConflictCached is optimized version using cached data
func checkConflictCached(cache *conflictCache, conflicts []Conflict, players []brackets.Player) (float64, int) {
	conflictScore := 0.0
	conflictSum := 0

	for _, s := range cache.conflictSets {
		for _, con := range conflicts {
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.checkCached(p1, p2, cache) {
				conflictScore += float64(2 + con.Priority)
				conflictSum++
			}
		}
	}

	return conflictScore, conflictSum
}

// listUnresolvedConflicts returns conflicts that were unresolved
func listUnresolvedConflicts(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) []Conflict {
	var unresolved []Conflict
	for _, s := range bracket.Sets {
		if s == nil {
			continue
		}
		for _, con := range conflicts {
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if con.check(p1, p2) {
				unresolved = append(unresolved, con)
			}
		}
	}
	return unresolved
}
