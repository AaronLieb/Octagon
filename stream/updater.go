// Package stream contains the core logic for automatically automatic updating the stream
package stream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/obs"
	"github.com/AaronLieb/octagon/parrygg"
	"github.com/AaronLieb/octagon/parrygg/pb"
	"github.com/charmbracelet/log"
)

func UpdateFromParryGG() error {
	apiKey := config.GetParryGGAPIKey()
	log.Infof("Parry.gg API Key: %s", apiKey)

	client, err := parrygg.NewClient(apiKey)
	if err != nil {
		return fmt.Errorf("failed to create parry.gg client: %w", err)
	}
	defer client.Close()

	ctx := client.WithAuth(context.Background())

	// Get cached stream match ID
	matchID, err := cache.GetStreamMatch()
	if len(matchID) <= 0 || err != nil {
		return fmt.Errorf("no stream match set: %w", err)
	}

	// Get bracket to access seeds and player info
	bracketResp, err := client.BracketService.GetBracket(ctx, &pb.GetBracketRequest{
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

	// Find the match in the bracket
	var match *pb.Match
	for _, m := range bracketResp.Bracket.Matches {
		if m.Id == matchID {
			match = m
			break
		}
	}

	if match == nil {
		return fmt.Errorf("stream match not found in bracket")
	}

	// Get player info from seeds
	var player1, player2 *pb.User
	if len(match.Slots) >= 2 {
		for _, seed := range bracketResp.Bracket.Seeds {
			if seed.Id == match.Slots[0].SeedId && seed.EventEntrant != nil && len(seed.EventEntrant.Entrant.Users) > 0 {
				player1 = seed.EventEntrant.Entrant.Users[0]
			}
			if seed.Id == match.Slots[1].SeedId && seed.EventEntrant != nil && len(seed.EventEntrant.Entrant.Users) > 0 {
				player2 = seed.EventEntrant.Entrant.Users[0]
			}
		}
	}

	if player1 == nil || player2 == nil {
		return fmt.Errorf("could not get player info from match")
	}

	log.Infof("Found match: %s - %s vs %s (State: %v)", match.Identifier, player1.GamerTag, player2.GamerTag, match.State)

	// Update OBS with current match state
	if err := updateOBSFromMatch(match, player1, player2); err != nil {
		return fmt.Errorf("failed to update OBS: %w", err)
	}

	log.Info("Successfully updated OBS from Parry.gg")
	return nil
}

func updateOBSFromMatch(match *pb.Match, player1, player2 *pb.User) error {
	// Get scores from match slots
	score1 := match.Slots[0].Score
	score2 := match.Slots[1].Score

	player1Score := strconv.FormatFloat(score1, 'f', 0, 64)
	player2Score := strconv.FormatFloat(score2, 'f', 0, 64)

	log.Infof("Updating OBS - %s (%s) vs %s (%s)", player1.GamerTag, player1Score, player2.GamerTag, player2Score)

	// Update OBS text sources
	if err := obs.UpdateText("Player 1 Name", player1.GamerTag); err != nil {
		return err
	}
	if err := obs.UpdateText("Player 2 Name", player2.GamerTag); err != nil {
		return err
	}
	if err := obs.UpdateText("Player 1 Score", player1Score); err != nil {
		return err
	}
	if err := obs.UpdateText("Player 2 Score", player2Score); err != nil {
		return err
	}

	return nil
}

func updateMatchScore(client *parrygg.Client, ctx context.Context, match *pb.Match, player1, player2 *pb.User, score1, score2 float64) error {
	log.Debugf("Creating match game for match %s with score %.0f-%.0f", match.Id, score1, score2)
	placement1 := int32(1)
	placement2 := int32(2)
	slotState := pb.SlotState_SLOT_STATE_NUMERIC

	_, err := client.MatchGameService.UpdateMatchGame(ctx, &pb.UpdateMatchGameRequest{
		Id: match.MatchGames[0].Id,
		MatchGame: &pb.MatchGameMutation{
			Index: 0,
			Slots: []*pb.GameSlotMutation{
				{
					Slot:      0,
					Score:     &score1,
					Placement: &placement1,
					Participants: []*pb.GameParticipantMutation{
						{UserId: player1.Id},
					},
					State: &slotState,
				},
				{
					Slot:      1,
					Score:     &score2,
					Placement: &placement2,
					Participants: []*pb.GameParticipantMutation{
						{UserId: player2.Id},
					},
					State: &slotState,
				},
			},
			State: pb.MatchGameState_MATCH_GAME_STATE_COMPLETED,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create match game: %w", err)
	}

	_, err = client.MatchGameService.CreateMatchGame(ctx, &pb.CreateMatchGameRequest{
		MatchId: match.Id,
		MatchGame: &pb.MatchGameMutation{
			Index: 1,
			Slots: []*pb.GameSlotMutation{
				{
					Slot:      0,
					Score:     &score1,
					Placement: &placement1,
					Participants: []*pb.GameParticipantMutation{
						{UserId: player1.Id},
					},
					State: &slotState,
				},
				{
					Slot:      1,
					Score:     &score2,
					Placement: &placement2,
					Participants: []*pb.GameParticipantMutation{
						{UserId: player2.Id},
					},
					State: &slotState,
				},
			},
			State: pb.MatchGameState_MATCH_GAME_STATE_COMPLETED,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create match game: %w", err)
	}

	return nil
}
