package attendees

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List attendees",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
		},
		Action: listAttendees,
	}
}

func listAttendees(ctx context.Context, cmd *cli.Command) error {
	log.Debug("list attendees")

	tournamentShortSlug := cmd.String("tournament")

	tournamentResp, err := startgg.GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}

	tournamentSlug := tournamentResp.Tournament.Slug

	// eventSlug := fmt.Sprintf("tournament/%s/event/ultimate-singles", tournamentSlug)

	resp, err := startgg.GetParticipants(ctx, tournamentSlug)
	if err != nil {
		return err
	}
	tournament := resp.Tournament
	participants := tournament.Participants.GetNodes()

	for _, participant := range participants {
		fmt.Printf("%s\t%s\t%s\n", participant.GamerTag, participant.ContactInfo.NameFirst, participant.ContactInfo.NameLast)
	}
	return nil
}
