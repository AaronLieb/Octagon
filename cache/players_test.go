package cache

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/AaronLieb/octagon/startgg"
)

func TestGetAllCachedPlayers(t *testing.T) {
	// Clear cache first
	Clear()
	
	// Add test players
	players := []startgg.CachedPlayer{
		{Name: "TestPlayer1", ID: "123"},
		{Name: "TestPlayer2", ID: "456"},
	}
	
	for _, player := range players {
		data, _ := json.Marshal(player)
		cacheKey := []byte("player_name:" + strings.ToLower(player.Name))
		Set(cacheKey, data)
	}
	
	// Test retrieval
	cached, err := GetAllCachedPlayers()
	if err != nil {
		t.Fatalf("GetAllCachedPlayers failed: %v", err)
	}
	
	if len(cached) != 2 {
		t.Errorf("Expected 2 players, got %d", len(cached))
	}
	
	// Verify players are correct
	found := make(map[string]bool)
	for _, player := range cached {
		found[player.Name] = true
	}
	
	if !found["TestPlayer1"] || !found["TestPlayer2"] {
		t.Error("Expected players not found in cache")
	}
}

func TestGetAllCachedPlayersEmpty(t *testing.T) {
	Clear()
	
	cached, err := GetAllCachedPlayers()
	if err != nil {
		t.Fatalf("GetAllCachedPlayers failed: %v", err)
	}
	
	if len(cached) != 0 {
		t.Errorf("Expected 0 players, got %d", len(cached))
	}
}

func TestPlayerCacheKeyFormat(t *testing.T) {
	player := startgg.CachedPlayer{Name: "TestPlayer", ID: "789"}
	data, _ := json.Marshal(player)
	
	// Test that keys are lowercase
	cacheKey := []byte("player_name:" + strings.ToLower(player.Name))
	Set(cacheKey, data)
	
	// Should be able to retrieve with exact key
	retrieved, err := Get(cacheKey)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	
	var cachedPlayer startgg.CachedPlayer
	err = json.Unmarshal(retrieved, &cachedPlayer)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	
	if cachedPlayer.Name != player.Name || cachedPlayer.ID != player.ID {
		t.Error("Cached player data doesn't match original")
	}
}
