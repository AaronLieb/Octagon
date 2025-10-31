package bracket

import (
	"context"
	"fmt"
	"strings"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const NameLength = 15

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

	b := brackets.CreateBracket(8)
	printRounds(b.WinnersRounds)
	// printRounds(b.LosersRounds)
	return nil
}

func printRounds(rounds [][]*brackets.Set) {
	height := len(rounds[0]) * 2
	width := NameLength * len(rounds)
	output := make([][]byte, height)
	for i := range output {
		output[i] = make([]byte, width)
		for j := range output[i] {
			output[i][j] = 'x'
		}
	}
	for y, round := range rounds {
		for i, set := range round {
			spacing := height / len(round)
			k := spacing/2 - 1
			Xoffset := NameLength * y
			Yoffset := k + i*spacing
			// probably need to use string builder
			p1Name := fmt.Sprintf("%2d┐", set.Player1)
			p2Name := fmt.Sprintf("%2d┘", set.Player2)
			log.Debug("name", "len", len(p1Name))
			p1Line := fmt.Sprintf("%s%s", p1Name, strings.Repeat("─", NameLength-len(p1Name)))
			p2Line := fmt.Sprintf("%s%s", p2Name, strings.Repeat("─", NameLength-len(p1Name)))
			// Why is this 35???
			log.Debug("line", "len", len(p1Line))
			log.Debug(Xoffset)
			for j := range NameLength {
				output[Yoffset][Xoffset+j] = p1Line[j]
				output[Yoffset+1][Xoffset+j] = p2Line[j]
			}
		}
	}
	for _, line := range output {
		fmt.Println(string(line))
	}
}
