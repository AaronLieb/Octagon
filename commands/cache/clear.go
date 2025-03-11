package cache

import (
	"context"

	"github.com/AaronLieb/octagon/cache"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func ClearCommand() *cli.Command {
	return &cli.Command{
		Name:    "clear",
		Usage:   "clears the cache",
		Aliases: []string{"c"},
		Action:  clear,
	}
}

func clear(ctx context.Context, cmd *cli.Command) error {
	err := cache.Clear()
	if err != nil {
		log.Fatalf("Unable to clear cache: %v", err)
	}
	return nil
}
