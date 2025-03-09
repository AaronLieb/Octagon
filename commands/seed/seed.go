package seed

import (
	"context"
	"errors"
	"math"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "seed",
		Usage:   "Seed a bracket",
		Aliases: []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "The slug for the tournament. Example: octagon-99",
				Value:   "octagon",
			},
			&cli.BoolFlag{
				Name:    "redemption",
				Aliases: []string{"r"},
				Usage:   "Whether you are seeding for redemption bracket or not",
			},
		},
		Action: seed,
	}
}

func seed(ctx context.Context, cmd *cli.Command) error {
	log.Debug("seed")

	if cmd.Args().Len() < 1 {
		return errors.New("missing argument: [numPlayers]")
	}

	numPlayers, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return err
	}

	createBracket(numPlayers)
	return nil
}

type set struct {
	name      string
	player1   int
	player2   int
	winnerSet *set
	loserSet  *set
}

var firstSets map[int]*set

func checkConflict(p1 int, p2 int) {
	// Checks the path of a given two players and returns whether they are seeded to play or not
}

/*
* This is a recursive function that generates sets and links them.
* The concept of a "round" is not your traditional bracket round.
* When referring to a seed as "high", the seed value is numerically low
* When referring to a seed that is "low" the seed value is numerically high
 */
func createSet(totRounds int, numPlayers int, round int, isWinners bool, highSeed int, winnerSet *set, loserSet *set) *set {
	numInRound := int(math.Pow(2, float64(totRounds-round)+1))

	lowSeed := numInRound - highSeed + 1

	/* This math is not simplified,
	 * but it makes sense in my head so
	 * I'm gonna keep it this way */
	if !isWinners {
		if highSeed <= numInRound {
			// second stage of losers round, middle seeds
			lb := numInRound/2 + 1
			ub := 2*numInRound - numInRound/2
			diff := highSeed - lb
			lowSeed = ub - diff

			// offset seeds so players don't play same as winners
			if diff >= numInRound/2 {
				lowSeed += numInRound / 4
			} else {
				lowSeed -= numInRound / 4
			}

		} else {
			// first stage of losers round, bottom seeds
			lb := numInRound + 1
			ub := numInRound * 2
			diff := highSeed - lb
			lowSeed = ub - diff
		}
	}

	// base case
	if lowSeed > numPlayers {
		return nil
	}

	newSet := &set{
		player1:   highSeed,
		player2:   lowSeed,
		winnerSet: winnerSet,
		loserSet:  loserSet,
	}
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", newSet.player1, newSet.player2, round, numInRound, isWinners)

	if !isWinners {
		return newSet
	}

	loserSet2 := createSet(totRounds, numPlayers, round, false, lowSeed, loserSet, nil)
	if loserSet2 != nil {
		loserSet1 := createSet(totRounds, numPlayers, round, false, loserSet2.player2, loserSet2, nil)
		if loserSet1 != nil {
			createSet(totRounds, numPlayers, round-1, true, highSeed, newSet, loserSet1)
			createSet(totRounds, numPlayers, round-1, true, lowSeed, newSet, loserSet1)
		}
	}

	return newSet
}

func createBracket(numPlayers int) *set {
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
	log.Debugf("%2d,%2d %2d %2d %t", winnersFinals.player1, winnersFinals.player2, numRounds, 2, true)
	log.Debugf("%2d,%2d %2d %2d %t", losersFinals.player1, losersFinals.player2, numRounds, 2, false)

	// Losers semis
	losersSemis := createSet(numRounds, numPlayers, numRounds, false, 3, losersFinals, nil)

	// Winners semis
	createSet(numRounds, numPlayers, numRounds-1, true, 1, winnersFinals, losersSemis)
	createSet(numRounds, numPlayers, numRounds-1, true, 2, winnersFinals, losersSemis)

	return winnersFinals
}
