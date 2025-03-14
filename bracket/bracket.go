package bracket

import (
	"math"

	"github.com/charmbracelet/log"
)

func CreateBracket(numPlayers int) Bracket {
	log.Debug("Create bracket", "players", numPlayers)

	numRounds := int(math.Ceil(math.Log2(float64(numPlayers))))

	losersFinals := &Set{
		Player1:   2,
		Player2:   3,
		WinnerSet: nil,
		LoserSet:  nil,
	}

	winnersFinals := &Set{
		Player1:   1,
		Player2:   2,
		WinnerSet: nil,
		LoserSet:  nil,
	}
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", winnersFinals.Player1, winnersFinals.Player2, numRounds, 2, true)
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", losersFinals.Player1, losersFinals.Player2, numRounds-1, 2, false)

	numSetsUpperBound := int(math.Pow(2, float64(numRounds+1)))
	sets := make([]*Set, numSetsUpperBound)

	// Losers semis
	losersSemis := createSet(numRounds, numPlayers, numRounds-1, false, 3, losersFinals, nil, &sets)

	// Winners semis
	createSet(numRounds, numPlayers, numRounds-2, true, 1, winnersFinals, losersSemis, &sets)
	createSet(numRounds, numPlayers, numRounds-2, true, 2, winnersFinals, losersSemis, &sets)

	var playableSets []*Set
	for _, set := range sets {
		if set != nil && set.Player1 <= numPlayers && set.Player2 <= numPlayers {
			playableSets = append(playableSets, set)
		}
	}

	log.Debug("playableSets", "len", len(playableSets))
	return Bracket{Sets: playableSets}
}

/*
* This is a recursive function that generates sets and links them.
* The concept of a "round" is not your traditional bracket round.
* When referring to a seed as "high", the seed value is numerically low
* When referring to a seed that is "low" the seed value is numerically high
 */
func createSet(totRounds int, numPlayers int, round int, isWinners bool, highSeed int, winnerSet *Set, loserSet *Set, sets *[]*Set) *Set {
	// base case
	if round < 1 {
		return nil
	}

	numInRound := int(math.Pow(2, float64(totRounds-round)))

	lowSeed := numInRound - highSeed + 1

	/* This math is not simplified,
	 * but it makes sense in my head so
	 * I'm gonna keep it this way */
	if !isWinners {
		if highSeed <= numInRound {
			// second stage of losers round, middle seeds
			lb := numInRound/2 + 1
			ub := lb + numInRound - 1
			diff := highSeed - lb
			lowSeed = ub - diff
			// log.Debug("", "lb", lb, "ub", ub, "diff", diff, "lowSeed", lowSeed)

			// offset seeds so players don't play same as winners
			dx := (round)%2 + 1
			if diff >= numInRound/2 {
				lowSeed -= dx
			} else {
				lowSeed += dx
			}
			// log.Debug("", "dx", dx, "lowSeed", lowSeed)

		} else {
			// first stage of losers round, bottom seeds
			lb := numInRound + 1
			ub := numInRound * 2
			diff := highSeed - lb
			lowSeed = ub - diff
		}
	}

	newSet := &Set{
		Player1:   highSeed,
		Player2:   lowSeed,
		WinnerSet: winnerSet,
		LoserSet:  loserSet,
	}

	// TODO: fix this shit
	x := 0
	for i, val := range *sets {
		if val == nil {
			x = i
		}
	}
	(*sets)[x] = newSet

	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", newSet.Player1, newSet.Player2, round, numInRound, isWinners)

	if !isWinners {
		return newSet
	}

	loserSet2 := createSet(totRounds, numPlayers, round, false, lowSeed, loserSet, nil, sets)
	loserSet1 := createSet(totRounds, numPlayers, round, false, loserSet2.Player2, loserSet2, nil, sets)
	createSet(totRounds, numPlayers, round-1, true, highSeed, newSet, loserSet1, sets)
	createSet(totRounds, numPlayers, round-1, true, lowSeed, newSet, loserSet1, sets)

	return newSet
}
