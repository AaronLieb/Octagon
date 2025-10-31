package conflicts

import (
	"time"

	"github.com/AaronLieb/octagon/startgg"
)

type Player struct {
	Name string     `json:"name"`
	ID   startgg.ID `json:"id"`
}

type Conflict struct {
	Priority   int        `json:"priority"`
	Reason     string     `json:"reason"`
	Players    []Player   `json:"players"`
	Expiration *time.Time `json:"expiration,omitempty"`
}
