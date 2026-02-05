package report

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/AaronLieb/octagon/tournament"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "report",
		Usage:   "Report tournament sets interactively",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "Tournament slug (default: octagon)",
				Value:   "octagon",
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage:   "Event name",
				Value:   "ultimate-singles",
			},
			&cli.BoolFlag{
				Name:    "redemption",
				Aliases: []string{"r"},
				Usage:   "Use redemption bracket",
			},
		},
		Action: runReportTUI,
	}
}

func runReportTUI(ctx context.Context, cmd *cli.Command) error {
	tournamentSlug, err := startgg.GetTournamentSlug(ctx, cmd.String("tournament"))
	if err != nil {
		return fmt.Errorf("tournament not found: %w", err)
	}

	event := cmd.String("event")
	if cmd.Bool("redemption") {
		event = startgg.EventRedemptionBracket
	}

	eventSlug := fmt.Sprintf(startgg.EventSlugFormat, tournamentSlug, event)

	sets, err := tournament.FetchReportableSets(ctx, eventSlug, false)
	if err != nil {
		return err
	}

	if len(sets) == 0 {
		log.Info("No reportable sets found")
		return nil
	}

	model := NewModel(ctx, sets)
	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
