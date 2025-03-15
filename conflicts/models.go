package conflicts

import "github.com/AaronLieb/octagon/startgg"

type player struct {
	Name string     `json:"name"`
	Id   startgg.ID `json:"id"`
}

type conflict struct {
	Priority int      `json:"priority"`
	Reason   string   `json:"reason"`
	Players  []player `json:"players"`
}
