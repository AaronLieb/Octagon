// Package brackets handles the logic for double elimination brackets, and determining which sets will be played
package brackets

import (
	"math"
)

/* Determined the flip and swap values by creating
 * brackets on startgg and manually checking
*
* "flipping" is reversing the order of the carryDown from winners
* Instead of the loser from the top of winners going to the top of losers,
* they will go to the bottom
*
* 5, 6, 7, 8 -> 8, 7, 6, 5
*
*
* "swapping" is splitting the winners in half, and swapping the groups
*
* 5, 6, 7, 8 -> 7, 8, 5, 6
*
* Flipping and swapping can be combined to get a "flip swap"
* 5, 6, 7, 8 -> 6, 5, 8, 7
* */

var flipMap map[int][]bool = map[int][]bool{
	8:   {true},
	16:  {true, false},
	32:  {true, true, true},
	64:  {true, true, false, false},
	128: {true, true, false, false, true},
}

var swapMap map[int][]bool = map[int][]bool{
	8:   {false},
	16:  {false, false},
	32:  {false, true, false},
	64:  {false, true, true, false},
	128: {false, true, true, false, false},
}

func createSets(round []int) []*Set {
	sets := make([]*Set, len(round)/2)

	for i := range sets {
		sets[i] = &Set{
			Player1: round[2*i],
			Player2: round[2*i+1],
		}
	}

	return sets
}

/* carryDown creates the "carryDown" loser round, which is when
* the losers from winners side play their first losers
* side set. The logic for this is strange, and it depends
* on the size of the bracket, and the round number.
 */
func carryDown(lr []int, wr []int, size int, round int) []int {
	r := make([]int, len(lr))

	flip := len(lr) > 2 && flipMap[size][round]
	swap := len(lr) > 2 && swapMap[size][round]

	wrl := reduceLosers(reduceWinners(wr))
	lrw := reduceWinners(lr)
	n := len(wrl)

	if flip {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			wrl[i], wrl[j] = wrl[j], wrl[i]
		}
	}

	if swap {
		wrl = append(wrl[n/2:n], wrl[0:n/2]...)
	}

	for i := range wrl {
		r[2*i] = wrl[i]
		r[2*i+1] = lrw[i]
	}
	return r
}

/* reduceWinners returns a new list with only
* the better seeds from each set. */
func reduceWinners(round []int) []int {
	half := make([]int, len(round)/2)
	for i := range half {
		half[i] = min(round[2*i], round[2*i+1])
	}
	return half
}

/* reduceLosers returns a new list with only
* the worse seeds from each set. */
func reduceLosers(round []int) []int {
	half := make([]int, len(round)/2)
	for i := range half {
		half[i] = max(round[2*i+1], round[2*i])
	}
	return half
}

func CreateBracket(numPlayers int) *Bracket {
	// the next power of 2, 36 person bracket -> n = 64
	n := int(math.Pow(2, math.Ceil(math.Log2(float64(numPlayers)))))

	var sets []*Set
	var losersRounds [][]*Set
	var winnersRounds [][]*Set

	// Initial winners bracket
	wr := CreateRound(n, 0)
	wrSets := createSets(wr)
	sets = append(sets, wrSets...)
	winnersRounds = append(winnersRounds, wrSets)

	// Initial losers bracket (byes)
	lr := CreateRound(n/2, n/2)
	lrSets := createSets(lr)
	sets = append(sets, lrSets...)
	losersRounds = append(losersRounds, lrSets)

	// Merge first winners bracket losers with initial losers
	lr = carryDown(lr, wr, n, 0)
	lrSets = createSets(lr)
	sets = append(sets, lrSets...)
	losersRounds = append(losersRounds, lrSets)

	for round := 1; len(wr) > 4; round++ {
		// Advance winners bracket
		wr = reduceWinners(wr)
		wrSets = createSets(wr)
		sets = append(sets, wrSets...)
		winnersRounds = append(winnersRounds, wrSets)

		// Advance losers bracket
		lr = reduceWinners(lr)
		lrSets = createSets(lr)
		sets = append(sets, lrSets...)
		losersRounds = append(losersRounds, lrSets)

		// Merge winners bracket losers with losers bracket
		lr = carryDown(lr, wr, n, round)
		lrSets = createSets(lr)
		sets = append(sets, lrSets...)
		losersRounds = append(losersRounds, lrSets)
	}

	// Winners bracket finals
	wr = reduceWinners(wr)
	wrSets = createSets(wr)
	sets = append(sets, wrSets...)
	winnersRounds = append(winnersRounds, wrSets)

	// Losers bracket finals
	for len(lr) > 2 {
		lr = reduceWinners(lr)
		lrSets = createSets(lr)
		sets = append(sets, lrSets...)
		losersRounds = append(losersRounds, lrSets)
	}

	// Filter out invalid players
	var setsFiltered []*Set
	for _, set := range sets {
		if set.Player1 <= numPlayers && set.Player2 <= numPlayers {
			setsFiltered = append(setsFiltered, set)
		}
	}

	bracket := &Bracket{
		Sets:          setsFiltered,
		WinnersRounds: winnersRounds,
		LosersRounds:  losersRounds,
	}

	return bracket
}

func CreateRound(n int, k int) []int {
	if n < 2 {
		return []int{}
	}

	if n == 2 {
		return []int{1, 2}
	}

	half := CreateRound(n/2, 0)
	newHalf := make([]int, len(half)*2)

	for i, num := range half {
		newHalf[2*i] = num + k
		newHalf[2*i+1] = n + 1 - num + k
	}

	return newHalf
}
