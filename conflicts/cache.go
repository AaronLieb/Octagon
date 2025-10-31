package conflicts

import (
	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/startgg"
)

// conflictCache holds pre-computed data for faster conflict resolution
type conflictCache struct {
	playerIndexMap map[startgg.ID]int
	conflictSets   []*brackets.Set
	stringCache    map[startgg.ID]string
}

func newConflictCache(players []brackets.Player, bracket *brackets.Bracket) *conflictCache {
	cache := &conflictCache{
		playerIndexMap: make(map[startgg.ID]int, len(players)),
		stringCache:    make(map[startgg.ID]string, len(players)),
	}
	
	// Pre-compute player index lookups
	for i, player := range players {
		cache.playerIndexMap[player.Id] = i
		cache.stringCache[player.Id] = startgg.ToString(player.Id)
	}
	
	// Pre-filter sets that could have conflicts
	for _, set := range bracket.Sets {
		if set != nil {
			cache.conflictSets = append(cache.conflictSets, set)
		}
	}
	
	return cache
}
