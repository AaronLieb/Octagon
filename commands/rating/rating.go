package rating

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/ratings"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "rating",
		Usage:  "get a player's rating",
		Action: GetRating,
	}
}

func GetRating(ctx context.Context, cmd *cli.Command) error {
	log.Debug("Rating")

	if cmd.Args().Len() < 1 {
		return fmt.Errorf("please provide a userid")
	}
	userId := cmd.Args().First()

	ratings.Get(ctx, userId)
	return nil
}
