package conflictscmd

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/conflicts"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "lists all conflicts",
		Aliases: []string{"l"},
		Action:  listConflict,
	}
}

func listConflict(ctx context.Context, cmd *cli.Command) error {
	cons := conflicts.GetConflicts([]string{})
	log.Info("found conflicts", "n", len(cons))
	for _, con := range cons {
		fmt.Printf("conflict p=%d\n", con.Priority)
		for _, p := range con.Players {
			fmt.Printf("  %s\n", p.Name)
		}
	}

	return nil
}
