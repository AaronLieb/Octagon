package stream

import (
	"context"
	"fmt"
	"time"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v3"
)

func WatchCommand() *cli.Command {
	return &cli.Command{
		Name:      "watch",
		Aliases:   []string{"w"},
		Usage:     "Listens for stream set updates from smartcv",
		UsageText: "octagon stream watch --port <port>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "Port to connect to",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "Tournament slug",
				Value:   "octagon",
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage:   "Event name",
				Value:   "ultimate-singles",
			},
		},
		Action: connectClient,
	}
}

func connectClient(ctx context.Context, cmd *cli.Command) error {
	port := cmd.Int("port")
	tournament := cmd.String("tournament")

	url := fmt.Sprintf("ws://localhost:%d/", port)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	fmt.Printf("Connected to %s\n", url)

	// Start goroutine to fetch stream queue every 30 seconds
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fetchStreamQueue(ctx, tournament)
			}
		}
	}()

	// Initial fetch
	fetchStreamQueue(ctx, tournament)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
	}

	err = conn.Close()
	if err != nil {
		log.Warn("Failed to close websocket connection", err)
	}

	return nil
}

func fetchStreamQueue(ctx context.Context, tournamentSlug string) {
	resp, err := startgg.GetStreamQueue(ctx, tournamentSlug)
	if err != nil {
		log.Printf("Failed to fetch stream queue: %v", err)
		return
	}

	if len(resp.Tournament.StreamQueue) == 0 {
		fmt.Println("Stream queue: empty")
		return
	}

	fmt.Printf("Stream queue (%s):\n", time.Now().Format("15:04:05"))
	for _, queue := range resp.Tournament.StreamQueue {
		streamName := queue.Stream.StreamName
		if streamName == "" {
			streamName = "Unknown"
		}

		// Find the latest started but incomplete set (state 2 = in progress)
		var currentSet *startgg.GetStreamQueueTournamentStreamQueueSetsSet
		for i := range queue.Sets {
			set := &queue.Sets[i]
			if set.State == 2 { // In progress
				currentSet = set
				break // Take the first in-progress set (should be latest)
			}
		}

		if currentSet == nil {
			fmt.Printf("  Stream: %s - No active sets\n", streamName)
			continue
		}

		if len(currentSet.Slots) >= 2 {
			p1 := "TBD"
			p2 := "TBD"

			if len(currentSet.Slots[0].Entrant.Participants) > 0 {
				p1 = currentSet.Slots[0].Entrant.Participants[0].Player.GamerTag
			}
			if len(currentSet.Slots[1].Entrant.Participants) > 0 {
				p2 = currentSet.Slots[1].Entrant.Participants[0].Player.GamerTag
			}

			round := currentSet.FullRoundText
			if round == "" {
				round = fmt.Sprintf("R%d", currentSet.Round)
			}

			fmt.Printf("  Stream: %s - %s: %s vs %s\n", streamName, round, p1, p2)
		}
	}
}
