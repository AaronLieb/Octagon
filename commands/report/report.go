package report

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

/*
* Command Examples
*
* 3-0 in Dahveed's favor King K Rool vs Wii Fit Trainer
* octagon report 3-0 dahveed Krool,wft
*
* 3-2 in Dahveed's favor Corrin vs Wii Fit Trainer
* octagon report 3-2 dahveed Wft,corrin Wft,corrin wft,Corrin wft,Corrin wft,Corrin
*
*
* Alternatively, we can make this an interactive command (maybe using gum)
* List off all the sets in progress, and have the user make a selection
* Ask for characters, and do it similarly to startgg
 */
func Command() *cli.Command {
	return &cli.Command{
		Name:    "report",
		Usage:   "Report a set",
		Aliases: []string{"r"},
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
				Usage:   "Whether you are reporting for redemption bracket or not",
			},
		},
		Action: report,
	}
}

func report(ctx context.Context, cmd *cli.Command) error {
	log.Debug("report")

	if cmd.Args().Len() < 2 {
		return fmt.Errorf("must provide score and at least one player")
	}

	player := cmd.Args().Get(1)

	// TODO: Clean this up
	score := strings.Split(cmd.Args().Get(0), "-")
	score1, err := strconv.Atoi(score[0])
	if err != nil {
		return fmt.Errorf("unable to parse score '%s'", cmd.Args().Get(0))
	}
	score2, err := strconv.Atoi(score[1])
	if err != nil {
		return fmt.Errorf("unable to parse score '%s'", cmd.Args().Get(0))
	}
	fmt.Println(score1, score2)

	tournamentSlug := cmd.String("tournament")

	event := "ultimate-singles"
	if cmd.Bool("redemption") {
		event = "redemption-bracket"
	}

	eventSlug := fmt.Sprintf("tournament/%s/event/%s", tournamentSlug, event)

	// Find Entrant Id
	entrantResp, err := startgg.GetEntrantByName(ctx, eventSlug, player)
	if err != nil {
		return err
	}
	entrant := entrantResp.Event.Entrants.Nodes[0]

	// Find Set in Event
	setsResp, err := startgg.GetSetsForEntrant(ctx, eventSlug, entrant.Id)
	if err != nil {
		return err
	}
	set := setsResp.Event.Sets.Nodes[0]
	fmt.Println(set.Id, set.Slots)

	// TODO: Actually report the set
	// startgg.ReportSet(ctx, set.Id, entrant.Id, {})

	return nil
}
