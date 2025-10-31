package rating

import (
	"fmt"
	"strings"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/agnivade/levenshtein"
)

func findPlayerByName(name string) (*startgg.CachedPlayer, error) {
	// Try exact match first
	cacheKey := []byte("player_name:" + strings.ToLower(name))
	if cached, err := cache.Get(cacheKey); err == nil {
		var player startgg.CachedPlayer
		if err := startgg.UnmarshalJSON(cached, &player); err == nil {
			return &player, nil
		}
	}

	// Get all cached players for similarity search
	allPlayers, err := cache.GetAllCachedPlayers()
	if err != nil {
		return nil, fmt.Errorf("failed to get cached players: %v", err)
	}

	if len(allPlayers) == 0 {
		return nil, fmt.Errorf("no players in cache - run 'octagon cache populate' first")
	}

	// Find best match using Levenshtein distance
	type match struct {
		player   startgg.CachedPlayer
		distance int
	}

	var matches []match
	searchName := strings.ToLower(name)

	for _, player := range allPlayers {
		playerName := strings.ToLower(player.Name)
		distance := levenshtein.ComputeDistance(searchName, playerName)

		// Only consider matches within reasonable edit distance
		maxDistance := max(len(name)/2, 2)

		if distance <= maxDistance {
			matches = append(matches, match{player, distance})
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found for '%s'", name)
	}

	return &matches[0].player, nil
}
