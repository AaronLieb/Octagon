package seed

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"

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

	// numPlayers, err := strconv.Atoi(cmd.Args().First())
	// if err != nil {
	// 	return err
	// }

	test()
	return nil
}

type set struct {
	name      string
	player1   int
	player2   int
	winnerSet *set
	loserSet  *set
}

type conflict struct {
	priority int
	players  []int // ID
}

type bracket struct {
	sets []*set
}

// returns true if p1 and p2 are in the conflict
func (con *conflict) check(p1 int, p2 int) bool {
	flag := true
	for _, p := range con.players {
		if p == p1 || p == p2 {
			if flag {
				flag = false
			} else {
				return true
			}
		}
	}
	return false
}

var firstSets map[int]*set

/*
* This is a recursive function that generates sets and links them.
* The concept of a "" is not your traditional bracket round.
* When referring to a seed as "high", the seed value is numerically low
* When referring to a seed that is "low" the seed value is numerically high
 */
func createSet(totRounds int, numPlayers int, round int, isWinners bool, highSeed int, winnerSet *set, loserSet *set, sets *[]*set) *set {
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

	// TODO fix this shit
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
	if loserSet2 != nil {
		loserSet1 := createSet(totRounds, numPlayers, round, false, loserSet2.player2, loserSet2, nil, sets)
		if loserSet1 != nil {
			createSet(totRounds, numPlayers, round-1, true, highSeed, newSet, loserSet1, sets)
			createSet(totRounds, numPlayers, round-1, true, lowSeed, newSet, loserSet1, sets)
		}
	}

	return newSet
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
	log.Debugf("s%-2d vs s%-2d r%-2d %2dp w:%t", losersFinals.player1, losersFinals.player2, numRounds, 2, false)

	sets := make([]*set, 50)
	sets[0] = winnersFinals
	sets[1] = losersFinals

	// Losers semis
	losersSemis := createSet(numRounds, numPlayers, numRounds, false, 3, losersFinals, nil, &sets)

	// Winners semis
	createSet(numRounds, numPlayers, numRounds-1, true, 1, winnersFinals, losersSemis, &sets)
	createSet(numRounds, numPlayers, numRounds-1, true, 2, winnersFinals, losersSemis, &sets)

	return bracket{sets: sets}
}

/*
 * Returns the conflict sum, which matches the following
 * 3 * p1_conflicts + 2 * p2_conflicts + p3_conflicts
 */
func checkConflict(b bracket, cons []conflict, players []int) int {
	conflictSum := 0

	// BFS
	for _, s := range b.sets {
		for _, con := range cons {
			if s == nil {
				break
			}
			p1 := players[s.player1-1]
			p2 := players[s.player2-1]
			if con.check(p1, p2) {
				log.Debug("conflict", "con", con, "p1", p1, "p2", p2)
				conflictSum += (4 - con.priority)
			}
		}
	}

	return conflictSum
}

/*
* Randomly shifts seeding
* Doesn't impact seeds 1 and 2
 */
func randomize(players []int) ([]int, int) {
	sum := 0
	newPlayers := make([]int, len(players))
	copy(newPlayers, players)
	for j := range players[3:] {
		i := j + 3
		r := rand.Reader
		x, _ := rand.Int(r, big.NewInt(5))
		if x.Int64() == 1 {
			temp := newPlayers[i]
			newPlayers[i] = newPlayers[i-1]
			newPlayers[i-1] = temp
			sum += 1
		}
	}
	return newPlayers, sum
}

func test() {
	players := []int{
		101,
		102,
		103,
		104,
		105,
		106,
		107,
		108,
	}

	var best []int
	lowest := 999
	for range 100 {
		newPlayers, seedDiff := randomize(players)
		log.Debug("test", "players", newPlayers, "seedDiff", seedDiff)

		conflicts := []conflict{
			{
				priority: 1,
				players:  []int{101, 108},
			},
			{
				priority: 2,
				players:  []int{103, 104},
			},
			{
				priority: 3,
				players:  []int{105, 103},
			},
			{
				priority: 3,
				players:  []int{107, 108},
			},
			{
				priority: 3,
				players:  []int{105, 104},
			},
			{
				priority: 1,
				players:  []int{101, 107},
			},
		}
		bracket := createBracket(len(players))
		sum := checkConflict(bracket, conflicts, newPlayers)
		sum += seedDiff
		if sum < lowest {
			lowest = sum
			best = newPlayers
		}
	}
	fmt.Println(lowest, best)
}
