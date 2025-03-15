package bracket

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/bracket"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func printCommand() *cli.Command {
	return &cli.Command{
		Name:    "print",
		Usage:   "print bracket",
		Aliases: []string{"p"},
		Action:  printBracket,
	}
}

func printBracket(ctx context.Context, cmd *cli.Command) error {
	log.Debug("print bracket")

	b := bracket.CreateBracket(32)
	height := len(b.WinnersRounds[0]) * 2
	width := 15 * len(b.WinnersRounds)
	output := make([][]byte, height)
	for i := range output {
		output[i] = make([]byte, width)
		for j := range output[i] {
			output[i][j] = ' '
		}
	}
	for y, round := range b.WinnersRounds {
		for i, set := range round {
			spacing := height / len(round)
			k := spacing/2 - 1
			if y == 0 {
				k = 0
			}
			offset := k + i*spacing
			output[offset][5*y] = fmt.Sprintf("%d", set.Player1)[0]
			output[offset+1][5*y] = fmt.Sprintf("%d", set.Player2)[0]
		}
	}
	for _, line := range output {
		fmt.Println(string(line))
	}
	return nil
}
