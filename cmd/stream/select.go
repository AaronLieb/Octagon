package stream

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/parrygg"
	"github.com/AaronLieb/octagon/parrygg/pb"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func SelectCommand() *cli.Command {
	return &cli.Command{
		Name:    "select",
		Usage:   "Select the stream match",
		Aliases: []string{"s"},
		Action:  runSelectTUI,
	}
}

func runSelectTUI(ctx context.Context, cmd *cli.Command) error {
	apiKey := config.GetParryGGAPIKey()
	client, err := parrygg.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("failed to create parry.gg client: %w", err)
	}
	defer client.Close()

	authCtx := client.WithAuth(context.Background())

	bracketResp, err := client.BracketService.GetBracket(authCtx, &pb.GetBracketRequest{
		Identifier: &pb.GetBracketRequest_SlugId{
			SlugId: &pb.BracketSlugId{
				TournamentSlug: "octagon",
				EventSlug:      "ultimate-singles",
				PhaseSlug:      "main",
				BracketSlug:    "bracket",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get bracket: %w", err)
	}

	matches := filterMatches(bracketResp.Bracket)
	if len(matches) == 0 {
		log.Info("No ready or started matches found")
		return nil
	}

	model := NewSelectModel(matches, bracketResp.Bracket.Seeds)
	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

func filterMatches(bracket *pb.Bracket) []*pb.Match {
	var matches []*pb.Match
	for _, match := range bracket.Matches {
		if match.State == pb.MatchState_MATCH_STATE_READY || match.State == pb.MatchState_MATCH_STATE_IN_PROGRESS {
			matches = append(matches, match)
		}
	}
	return matches
}
