package conflicts

import (
	"context"

	"github.com/urfave/cli/v3"
)

func addCommand() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Usage:   "List attendees",
		Aliases: []string{"create", "c", "a"},
		Action:  addConflict,
	}
}

func addConflict(ctx context.Context, cmd *cli.Command) error {
	return nil
}
