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
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "The output type of the attendee list. (table, csv)",
				Value:   "table",
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

	resp, err := startgg.GetParticipants(ctx, tournamentSlug)
	if err != nil {
		return err
	}
	tournament := resp.Tournament
	participants := tournament.Participants.GetNodes()

	log.Infof("Fetched %d participants for %s", len(participants), tournamentSlug)

	outputType := cmd.String("output")
	if outputType == "csv" {
		for _, participant := range participants {
			fmt.Printf("%s\t%s\t%s\n", participant.GamerTag, participant.ContactInfo.NameFirst, participant.ContactInfo.NameLast)
		}
	} else {
		for _, participant := range participants {
			player := participant.Player
			fmt.Printf("%-25s %-8d %-15s %-15s\n", participant.GamerTag, player.Id, participant.ContactInfo.NameFirst, participant.ContactInfo.NameLast)
		}
	}
	return nil
}
