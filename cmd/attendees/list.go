package attendees

import (
	"context"
	"strings"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
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

	resp, err := startgg.GetParticipants(ctx, tournamentShortSlug)
	if err != nil {
		return err
	}
	tournament := resp.Tournament
	participants := tournament.Participants.GetNodes()

	log.Infof("Fetched %d participants for %s", len(participants), tournament.Name)

	outputType := strings.ToLower(cmd.String("output"))
	for _, participant := range participants {
		player := participant.Player
		if outputType == "csv" {
			startgg.PrintPlayerCSV(player.Id, participant.GamerTag)
		} else {
			startgg.PrintPlayerTable(participant.GamerTag, player.Id, participant.ContactInfo.NameFirst, participant.ContactInfo.NameLast)
		}
	}
	return nil
}
