package rating

import (
	"context"
	"fmt"
	"strconv"

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
	userIdStr := cmd.Args().First()

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return err
	}

	rating, err := ratings.Get(ctx, userId)
	if err != nil {
		return err
	}
	fmt.Println(rating)
	return nil
}
