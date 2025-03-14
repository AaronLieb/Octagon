package attendees

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func addCommand() *cli.Command {
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
			&cli.BoolFlag{
				Name:    "redemption",
				Aliases: []string{"r"},
				Usage:   "Add to redemption bracket",
			},
		},
		Action: addAttendee,
	}
}

/* This will not actually work, the registerForTournament
 * query has to be made by the user registering.
 */
func addAttendee(ctx context.Context, cmd *cli.Command) error {
	log.Debug("add attendee")

	log.Warn("This does not work and will not work until startgg updates their api")

	if cmd.Args().Len() < 1 {
		return fmt.Errorf("missing argument: [userId]")
	}

	userId, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return fmt.Errorf("unable to parse argument [userId]: %v", err)
	}

	tournamentShortSlug := cmd.String("tournament")

	tournamentResp, err := startgg.GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	tournamentSlug := tournamentResp.Tournament.Slug

	eventName := "ultimate-singles"
	if cmd.Bool("redemption") {
		eventName = "redemption-bracket"
	}

	eventSlug := fmt.Sprintf("%s/event/%s", tournamentSlug, eventName)
	eventResp, err := startgg.GetEvent(ctx, eventSlug)
	if err != nil {
		return err
	}
	event := eventResp.Event

	tokenResp, err := startgg.GenerateRegistrationToken(ctx, event.Id, userId)
	if err != nil {
		return err
	}
	regToken := tokenResp.GenerateRegistrationToken
	log.Debug("generate registration token", "token", regToken)

	registerResp, err := startgg.RegisterForTournament(ctx, event.Id, regToken)
	if err != nil {
		return err
	}
	log.Info("Successfully registered", "id", registerResp.RegisterForTournament.Id)

	return nil
}
