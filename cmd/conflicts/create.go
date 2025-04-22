package conflicts

import (
	"context"

	"github.com/urfave/cli/v3"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Usage:   "creates a conflict",
		Aliases: []string{"c", "add", "a"},
		Action:  createConflict,
	}
}

func createConflict(ctx context.Context, cmd *cli.Command) error {
	return nil
}
