package reportcmd

import (
	"context"
	"fmt"

	sgg "github.com/AaronLieb/octagon/commands/reportcmd/startgg"
	"github.com/AaronLieb/octagon/commands/reportcmd/ui"
	"github.com/AaronLieb/octagon/startgg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func testCommand() *cli.Command {
	return &cli.Command{
		Name:    "test",
		Usage:   "test command",
		Aliases: []string{"r"},
		Action:  test,
	}
}

func test(ctx context.Context, cmd *cli.Command) error {
	tournamentName := "octagon"
	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentName)
	if err != nil {
		log.Fatalf("unable to find tournament '%s': %v", tournamentName, err)
	}

	eventSlug := fmt.Sprintf("%s/event/ultimate-singles", tournamentSlug)
	sets := sgg.FetchReportableSets(ctx, eventSlug)

	p := tea.NewProgram(ui.InitialModel(sets))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
