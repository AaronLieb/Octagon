package attendees

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func AddCommand() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Usage:   "Add attendee",
		Aliases: []string{"a"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
		},
		Action: addAttendee,
	}
}

func addAttendee(ctx context.Context, cmd *cli.Command) error {
	log.Debug("add attendee")

	currentUserResp, err := startgg.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	currentUser := currentUserResp.CurrentUser

	fmt.Println(currentUser)

	reg, err := startgg.GenerateRegistrationToken(ctx, currentUser.Id)
	if err != nil {
		return err
	}

	fmt.Println(reg)
	return nil
}
