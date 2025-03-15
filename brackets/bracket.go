package brackets

import (
	"math"
)

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

func carryDown(lr []int, wr []int) []int {
	r := make([]int, len(lr))
	for i := 0; i < len(lr); i += 2 {
		r[i] = wr[len(wr)-i*2-2]
		r[i+1] = lr[i]
	}
	return r
}

func reduce(round []int) []int {
	half := make([]int, len(round)/2)
	for i := range half {
		half[i] = round[2*i]
	}
	return half
}

// temporary code that I need to refactor
// way too much reptition, logic can be pulled out
// to funcs
func CreateBracket(numPlayers int) *Bracket {
	n := int(math.Pow(2, math.Ceil(math.Log2(float64(numPlayers)))))

	var sets []*Set
	var losersRounds [][]*Set
	var winnersRounds [][]*Set
	var wrSets, lr1Sets, lr2Sets []*Set
	wr := CreateRound(n, 0)
	lr1 := CreateRound(n/2, n/2)
	lr2 := carryDown(lr1, wr)
	wrSets = createSets(wr)
	lr1Sets = createSets(lr1)
	lr2Sets = createSets(lr2)
	sets = append(sets, wrSets...)
	sets = append(sets, lr1Sets...)
	sets = append(sets, lr2Sets...)
	winnersRounds = append(winnersRounds, wrSets)
	losersRounds = append(losersRounds, lr1Sets)
	losersRounds = append(losersRounds, lr2Sets)

	for len(wr) > 4 {
		wr = reduce(wr)
		lr1 = reduce(lr2)
		lr2 = carryDown(lr1, wr)

		wrSets = createSets(wr)
		lr1Sets = createSets(lr1)
		lr2Sets = createSets(lr2)
		sets = append(sets, wrSets...)
		sets = append(sets, lr1Sets...)
		sets = append(sets, lr2Sets...)
		winnersRounds = append(winnersRounds, wrSets)
		losersRounds = append(losersRounds, lr1Sets)
		losersRounds = append(losersRounds, lr2Sets)
	}

	wr = reduce(wr)
	wrSets = createSets(wr)
	sets = append(sets, wrSets...)
	winnersRounds = append(winnersRounds, wrSets)

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
