package cache

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/urfave/cli/v3"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List cached players",
		Aliases: []string{"l", "ls"},
		Action:  listCache,
	}
}

func listCache(ctx context.Context, cmd *cli.Command) error {
	players, err := cache.GetAllCachedPlayers()
	if err != nil {
		return err
	}

	if len(players) == 0 {
		fmt.Println("No players in cache")
		return nil
	}

	fmt.Printf("Found %d cached players:\n\n", len(players))
	for _, player := range players {
		fmt.Printf("%-30s %s\n", player.Name, startgg.ToString(player.ID))
	}

	return nil
}
