package rating

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/urfave/cli/v3"
)

func biasCommand() *cli.Command {
	return &cli.Command{
		Name:  "bias",
		Usage: "Manage rating biases",
		Commands: []*cli.Command{
			addBiasCommand(),
			listBiasCommand(),
		},
	}
}

func addBiasCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a rating bias: octagon rating bias add playerId ratio reason",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "expires",
				Aliases: []string{"e"},
				Usage:   "Expiration duration (e.g., '24h', '7d', '2w')",
			},
		},
		Action: addBias,
	}
}

func listBiasCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all rating biases",
		Action: listBiases,
	}
}

func addBias(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("usage: octagon rating bias add playerId ratio reason")
	}

	playerInput := args[0]
	ratioStr := args[1]
	reason := strings.Join(args[2:], " ")

	// Try to find player by name, fall back to ID
	playerID, err := resolvePlayerInput(playerInput)
	if err != nil {
		return fmt.Errorf("could not resolve player '%s': %v", playerInput, err)
	}

	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil {
		return fmt.Errorf("invalid ratio: %v", err)
	}

	if ratio <= 0 {
		return fmt.Errorf("ratio must be positive")
	}

	var expiration *time.Time
	if expiresStr := cmd.String("expires"); expiresStr != "" {
		duration, err := parseDuration(expiresStr)
		if err != nil {
			return fmt.Errorf("invalid expiration format: %v", err)
		}
		exp := time.Now().Add(duration)
		expiration = &exp
	}

	bias := config.Bias{
		PlayerID:   playerID,
		Ratio:      ratio,
		Reason:     reason,
		Expiration: expiration,
	}

	err = config.SaveBias(bias)
	if err != nil {
		return fmt.Errorf("failed to save bias: %v", err)
	}

	fmt.Printf("Added bias for player %s: %.2fx (%s)\n", playerID, ratio, reason)
	if expiration != nil {
		fmt.Printf("Expires: %s\n", expiration.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func resolvePlayerInput(input string) (string, error) {
	// Try to find player by name first
	player, err := findPlayerByName(input)
	if err == nil {
		return startgg.ToString(player.ID), nil
	}

	// Fall back to treating input as ID
	if _, err := strconv.Atoi(input); err == nil {
		return input, nil
	}

	return "", fmt.Errorf("not found as player name or valid ID")
}

func listBiases(ctx context.Context, cmd *cli.Command) error {
	biases := config.GetBiases()

	if len(biases) == 0 {
		fmt.Println("No biases found")
		return nil
	}

	fmt.Printf("%-15s %-10s %-20s %s\n", "Player ID", "Ratio", "Expires", "Reason")
	fmt.Println(strings.Repeat("-", 70))

	for _, bias := range biases {
		expiresStr := "Never"
		if bias.Expiration != nil {
			expiresStr = bias.Expiration.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-15s %8.2fx %-20s %s\n",
			bias.PlayerID, bias.Ratio, expiresStr, bias.Reason)
	}

	return nil
}

func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		if d, err := time.ParseDuration(days + "h"); err == nil {
			return d * 24, nil
		}
	}

	if strings.HasSuffix(s, "w") {
		weeks := strings.TrimSuffix(s, "w")
		if d, err := time.ParseDuration(weeks + "h"); err == nil {
			return d * 168, nil
		}
	}

	return time.ParseDuration(s)
}
