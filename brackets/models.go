package brackets

import "github.com/AaronLieb/octagon/startgg"

type Player struct {
	Name   string
	ID     startgg.ID
	Rating float64
}

type Set struct {
	Name      string
	Player1   int // good seed
	Player2   int // bad seed
	WinnerSet *Set
	LoserSet  *Set
}

type Bracket struct {
	Sets          []*Set
	WinnersRounds [][]*Set
	LosersRounds  [][]*Set
}
