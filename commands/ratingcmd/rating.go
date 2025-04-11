package ratingcmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AaronLieb/octagon/ratings"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const desc = `Will fetch player ratings from an external rating database.
Currently, this command will not function without credentials to a
firebase realtime database containing player ratings.`

func Command() *cli.Command {
	return &cli.Command{
		Name:        "rating",
		Usage:       "Fetch a player's rating",
		Description: desc,
		UsageText:   "octagon rating [playerId]",
		Action:      GetRating,
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
