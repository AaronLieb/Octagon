package brackets

import (
	"testing"
)

func assertSetExistsInRound(t *testing.T, bracket *Bracket, s1 int, s2 int, roundNum int, isWinners bool) {
	var round []*Set
	if isWinners {
		round = bracket.WinnersRounds[roundNum]
	} else {
		round = bracket.LosersRounds[roundNum]
	}
	for _, set := range round {
		if set.Player1 == s1 && set.Player2 == s2 {
			return
		}
	}
	t.Fatalf(`Expected to find set {%d, %d}`, s1, s2)
}

func TestCreateBracket(t *testing.T) {
	b := CreateBracket(8)

	assertSetExistsInRound(t, b, 1, 8, 0, true)
	assertSetExistsInRound(t, b, 2, 7, 0, true)
	assertSetExistsInRound(t, b, 3, 6, 0, true)
}
