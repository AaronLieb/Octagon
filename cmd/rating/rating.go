// Package rating is responsibel for fetching and manipulating SchuStats ratings
package rating

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
		Commands: []*cli.Command{
			biasCommand(),
		},
	}
}

func GetRating(ctx context.Context, cmd *cli.Command) error {
	log.Debug("Rating")

	if cmd.Args().Len() < 1 {
		return fmt.Errorf("please provide a userid or player name")
	}
	userInput := cmd.Args().First()

	// Try to resolve as player name first, then as ID
	var userID int
	player, err := findPlayerByName(userInput)
	if err == nil {
		// Convert ID to int
		if id, ok := player.ID.(float64); ok {
			userID = int(id)
		} else if id, ok := player.ID.(int); ok {
			userID = id
		} else {
			return fmt.Errorf("invalid player ID type")
		}
	} else {
		// Try as direct ID
		userID, err = strconv.Atoi(userInput)
		if err != nil {
			return fmt.Errorf("could not resolve '%s' as player name or ID: %v", userInput, err)
		}
	}

	rating, err := ratings.Get(ctx, userID)
	if err != nil {
		return err
	}
	fmt.Println(rating)
	return nil
}
