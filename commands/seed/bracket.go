package seed

import (
	"math"

	"github.com/charmbracelet/log"
)

type set struct {
	name      string
	player1   int
	player2   int
	winnerSet *set
	loserSet  *set
}

type bracket struct {
	sets []*set
}

func createBracket(numPlayers int) bracket {
	log.Debug("Create bracket", "players", numPlayers)

	numRounds := int(math.Ceil(math.Log2(float64(numPlayers))))

	losersFinals := &set{
		player1:   2,
		player2:   3,
		winnerSet: nil,
		loserSet:  nil,
	}

	winnersFinals := &set{
		player1:   1,
		player2:   2,
		winnerSet: nil,
		loserSet:  nil,
	}
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", winnersFinals.player1, winnersFinals.player2, numRounds, 2, true)
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", losersFinals.player1, losersFinals.player2, numRounds-1, 2, false)

	numSetsUpperBound := int(math.Pow(2, float64(numRounds+1)))
	sets := make([]*set, numSetsUpperBound)

	// Losers semis
	losersSemis := createSet(numRounds, numPlayers, numRounds-1, false, 3, losersFinals, nil, &sets)

	// Winners semis
	createSet(numRounds, numPlayers, numRounds-2, true, 1, winnersFinals, losersSemis, &sets)
	createSet(numRounds, numPlayers, numRounds-2, true, 2, winnersFinals, losersSemis, &sets)

	var playableSets []*set
	for _, set := range sets {
		if set != nil && set.player1 <= numPlayers && set.player2 <= numPlayers {
			playableSets = append(playableSets, set)
		}
	}

	log.Debug("playableSets", "len", len(playableSets))
	return bracket{sets: playableSets}
}

/*
* This is a recursive function that generates sets and links them.
* The concept of a "" is not your traditional bracket round.
* When referring to a seed as "high", the seed value is numerically low
* When referring to a seed that is "low" the seed value is numerically high
 */
func createSet(totRounds int, numPlayers int, round int, isWinners bool, highSeed int, winnerSet *set, loserSet *set, sets *[]*set) *set {
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

	newSet := &set{
		player1:   highSeed,
		player2:   lowSeed,
		winnerSet: winnerSet,
		loserSet:  loserSet,
	}

	// TODO: fix this shit
	x := 0
	for i, val := range *sets {
		if val == nil {
			x = i
		}
	}
	(*sets)[x] = newSet

	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", newSet.player1, newSet.player2, round, numInRound, isWinners)

	if !isWinners {
		return newSet
	}

	loserSet2 := createSet(totRounds, numPlayers, round, false, lowSeed, loserSet, nil, sets)
	loserSet1 := createSet(totRounds, numPlayers, round, false, loserSet2.player2, loserSet2, nil, sets)
	createSet(totRounds, numPlayers, round-1, true, highSeed, newSet, loserSet1, sets)
	createSet(totRounds, numPlayers, round-1, true, lowSeed, newSet, loserSet1, sets)

	return newSet
}
