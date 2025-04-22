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

	// wr1
	assertSetExistsInRound(t, b, 1, 8, 0, true)
	assertSetExistsInRound(t, b, 2, 7, 0, true)
	assertSetExistsInRound(t, b, 3, 6, 0, true)

	// lr1
	assertSetExistsInRound(t, b, 5, 8, 0, false)
	assertSetExistsInRound(t, b, 6, 7, 0, false)

	// lr2
	assertSetExistsInRound(t, b, 4, 6, 1, false)
	assertSetExistsInRound(t, b, 3, 5, 1, false)

	// lr3
	assertSetExistsInRound(t, b, 3, 4, 2, false)

	// lr4
	assertSetExistsInRound(t, b, 2, 3, 3, false)

	b = CreateBracket(16)

	assertSetExistsInRound(t, b, 9, 16, 0, false)
	assertSetExistsInRound(t, b, 6, 9, 1, false)

	assertSetExistsInRound(t, b, 12, 13, 0, false)
	assertSetExistsInRound(t, b, 7, 12, 1, false)

	assertSetExistsInRound(t, b, 8, 11, 1, false)
	assertSetExistsInRound(t, b, 5, 8, 2, false)

	b = CreateBracket(32)

	// winners r1
	assertSetExistsInRound(t, b, 1, 32, 0, true)
	assertSetExistsInRound(t, b, 16, 17, 0, true)
	assertSetExistsInRound(t, b, 4, 29, 0, true)
	assertSetExistsInRound(t, b, 13, 20, 0, true)

	// losers r1
	assertSetExistsInRound(t, b, 21, 28, 0, false)
	assertSetExistsInRound(t, b, 20, 29, 0, false)
	assertSetExistsInRound(t, b, 22, 27, 0, false)
	assertSetExistsInRound(t, b, 18, 31, 0, false)

	// losers r2
	assertSetExistsInRound(t, b, 12, 18, 1, false)
	assertSetExistsInRound(t, b, 16, 22, 1, false)
	assertSetExistsInRound(t, b, 13, 23, 1, false)

	// losers r3
	assertSetExistsInRound(t, b, 12, 13, 2, false)
	assertSetExistsInRound(t, b, 9, 16, 2, false)

	b = CreateBracket(64)

	// winners r1
	assertSetExistsInRound(t, b, 22, 43, 0, true)
	assertSetExistsInRound(t, b, 11, 54, 0, true)
	assertSetExistsInRound(t, b, 27, 38, 0, true)
	assertSetExistsInRound(t, b, 27, 38, 0, true)
	assertSetExistsInRound(t, b, 6, 59, 0, true)
	assertSetExistsInRound(t, b, 30, 35, 0, true)
	assertSetExistsInRound(t, b, 23, 42, 0, true)
	assertSetExistsInRound(t, b, 7, 58, 0, true)
	assertSetExistsInRound(t, b, 2, 63, 0, true)

	// losers r1
	assertSetExistsInRound(t, b, 43, 54, 0, false)

	// losers r2
	assertSetExistsInRound(t, b, 32, 43, 1, false)
	assertSetExistsInRound(t, b, 20, 39, 1, false)
	assertSetExistsInRound(t, b, 22, 33, 1, false)

	// losers r3
	assertSetExistsInRound(t, b, 21, 28, 2, false)
	assertSetExistsInRound(t, b, 19, 30, 2, false)

	// losers r3
	assertSetExistsInRound(t, b, 11, 21, 3, false)
	assertSetExistsInRound(t, b, 14, 20, 3, false)
	assertSetExistsInRound(t, b, 9, 23, 3, false)
}
