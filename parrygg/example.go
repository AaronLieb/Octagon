package parrygg

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/parrygg/pb"
)

// Example usage of the Parry.gg client
func ExampleGetTournaments(apiKey string) error {
	client, err := NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	ctx := client.WithAuth(context.Background())
	
	req := &pb.GetTournamentsRequest{}
	resp, err := client.TournamentService.GetTournaments(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get tournaments: %w", err)
	}

	fmt.Printf("Found %d tournaments\n", len(resp.Tournaments))
	return nil
}
