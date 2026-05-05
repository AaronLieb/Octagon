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
	setRoundMap    map[*brackets.Set]brackets.Round
}

func newConflictCache(players []brackets.Player, bracket *brackets.Bracket) *conflictCache {
	cache := &conflictCache{
		playerIndexMap: make(map[startgg.ID]int, len(players)),
		stringCache:    make(map[startgg.ID]string, len(players)),
		setRoundMap:    make(map[*brackets.Set]brackets.Round),
	}

	// Pre-compute player index lookups
	for i, player := range players {
		cache.playerIndexMap[player.ID] = i
		cache.stringCache[player.ID] = startgg.ToString(player.ID)
	}

	// Pre-filter sets that could have conflicts
	for _, set := range bracket.Sets {
		if set != nil {
			cache.conflictSets = append(cache.conflictSets, set)
		}
	}

	// Map sets to their round (1-indexed to match ParseRound / Round.String /
	// RoundFromStartGG).
	for i, sets := range bracket.WinnersRounds {
		for _, set := range sets {
			cache.setRoundMap[set] = brackets.Round{Number: i + 1}
		}
	}
	for i, sets := range bracket.LosersRounds {
		for _, set := range sets {
			cache.setRoundMap[set] = brackets.Round{Losers: true, Number: i + 1}
		}
	}

	return cache
}
