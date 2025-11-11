package conflicts

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/conflicts"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/agnivade/levenshtein"
	"github.com/urfave/cli/v3"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Usage:   "creates a conflict: octagon conflict create player1name player2name",
		Aliases: []string{"c", "add", "a"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "expires",
				Aliases: []string{"e"},
				Usage:   "Expiration duration (e.g., '24h', '7d', '2w')",
			},
			&cli.IntFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "Priority level (1-3)",
				Value:   3,
			},
			&cli.StringFlag{
				Name:    "reason",
				Aliases: []string{"r"},
				Usage:   "Reason for the conflict",
			},
		},
		Action: createConflict,
	}
}

func createConflict(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("usage: octagon conflict create player1name player2name")
	}

	player1Name := args[0]
	player2Name := args[1]
	reason := cmd.String("reason")

	// Validate priority
	priority := cmd.Int("priority")
	if priority < 1 || priority > 3 {
		return fmt.Errorf("priority must be between 1 and 3")
	}

	// Parse expiration if provided
	var expiration *time.Time
	if expiresStr := cmd.String("expires"); expiresStr != "" {
		duration, err := parseDuration(expiresStr)
		if err != nil {
			return fmt.Errorf("invalid expiration format: %v", err)
		}
		exp := time.Now().Add(duration)
		expiration = &exp
	}

	// Find player IDs
	player1, err := findPlayer(player1Name)
	if err != nil {
		return fmt.Errorf("could not find player '%s': %v", player1Name, err)
	}

	player2, err := findPlayer(player2Name)
	if err != nil {
		return fmt.Errorf("could not find player '%s': %v", player2Name, err)
	}

	// Confirm players
	fmt.Printf("Create conflict between:\n")
	fmt.Printf("  Player 1: %s (ID: %s)\n", player1.Name, startgg.ToString(player1.ID))
	fmt.Printf("  Player 2: %s (ID: %s)\n", player2.Name, startgg.ToString(player2.ID))
	fmt.Printf("  Priority: %d\n", priority)
	if reason != "" {
		fmt.Printf("  Reason: %s\n", reason)
	}
	fmt.Print("Confirm? (y/N): ")

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("failed to parse input")
	}
	if input != "y" && input != "Y" {
		fmt.Println("Cancelled")
		return nil
	}

	// Create and save conflict
	newConflict := conflicts.Conflict{
		Priority:   int(priority),
		Reason:     reason,
		Expiration: expiration,
		Players: []conflicts.Player{
			{Name: player1.Name, ID: player1.ID},
			{Name: player2.Name, ID: player2.ID},
		},
	}

	err = conflicts.SaveConflict(newConflict)
	if err != nil {
		return fmt.Errorf("failed to save conflict: %v", err)
	}

	fmt.Println("Conflict created successfully")
	return nil
}

func findPlayer(name string) (*startgg.CachedPlayer, error) {
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

	// Sort by distance (best matches first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	return &matches[0].player, nil
}

func parseDuration(s string) (time.Duration, error) {
	// Handle days
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		if d, err := time.ParseDuration(days + "h"); err == nil {
			return d * 24, nil
		}
	}
	
	// Handle weeks  
	if strings.HasSuffix(s, "w") {
		weeks := strings.TrimSuffix(s, "w")
		if d, err := time.ParseDuration(weeks + "h"); err == nil {
			return d * 168, nil // 24 * 7
		}
	}
	
	// Standard Go duration parsing for h, m, s
	return time.ParseDuration(s)
}
