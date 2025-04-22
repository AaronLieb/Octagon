package brackets

import (
	"fmt"
	"testing"
)

func setExistsInRound(bracket *Bracket, s1 int, s2 int, roundNum int, isWinners bool) bool {
	var round []*Set
	if isWinners {
		round = bracket.WinnersRounds[roundNum]
	} else {
		round = bracket.LosersRounds[roundNum]
	}
	for _, set := range round {
		if set.Player1 == s1 && set.Player2 == s2 {
			return true
		}
	}
	return false
}

func TestCreateBracket(t *testing.T) {
	type testSet struct {
		p1      int
		p2      int
		round   int
		winners bool
	}
	tests := make(map[int][]testSet)
	tests[8] = []testSet{
		{1, 8, 0, true},
		{2, 7, 0, true},
		{3, 6, 0, true},
		{5, 8, 0, false},
		{6, 7, 0, false},

		{4, 6, 1, false},
		{3, 5, 1, false},

		{3, 4, 2, false},
		{2, 3, 3, false},
	}
	tests[16] = []testSet{
		{9, 16, 0, false},
		{6, 9, 1, false},
		{12, 13, 0, false},
		{7, 12, 1, false},
		{8, 11, 1, false},
		{5, 8, 2, false},
	}

	tests[32] = []testSet{
		{1, 32, 0, true},
		{16, 17, 0, true},
		{4, 29, 0, true},
		{13, 20, 0, true},

		{21, 28, 0, false},
		{20, 29, 0, false},
		{22, 27, 0, false},
		{18, 31, 0, false},

		{12, 18, 1, false},
		{16, 22, 1, false},
		{13, 23, 1, false},

		{12, 13, 2, false},
		{9, 16, 2, false},
	}

	tests[64] = []testSet{
		{22, 43, 0, true},
		{11, 54, 0, true},
		{6, 59, 0, true},

		{43, 54, 0, false},

		{32, 43, 1, false},
		{20, 39, 1, false},
		{22, 33, 1, false},

		{21, 28, 2, false},
		{19, 30, 2, false},

		{11, 21, 3, false},
		{14, 20, 3, false},
		{9, 23, 3, false},
	}

	for bracketSize, test := range tests {
		b := CreateBracket(bracketSize)
		for _, set := range test {
			t.Run(fmt.Sprintf("Bracket Size %d", bracketSize), func(t *testing.T) {
				if !setExistsInRound(b, set.p1, set.p2, set.round, set.winners) {
					t.Fatalf(`Expected to find set {%d, %d} r:%d, w:%t n:%d`, set.p1, set.p2, set.round, set.winners, bracketSize)
				}
			})
		}
	}
}
