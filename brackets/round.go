package brackets

import (
	"fmt"
	"strconv"
	"strings"
)

// Round represents a specific round in a double elimination bracket.
// Winners rounds are "WR1", "WR2", etc. Losers rounds are "LR1", "LR2", etc.
type Round struct {
	Losers bool
	Number int
}

// ParseRound parses a round string like "WR2" or "LR1" into a Round.
func ParseRound(s string) (Round, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	var losers bool
	var numStr string

	switch {
	case strings.HasPrefix(s, "WR"):
		numStr = s[2:]
	case strings.HasPrefix(s, "LR"):
		losers = true
		numStr = s[2:]
	default:
		return Round{}, fmt.Errorf("invalid round format %q: must start with WR or LR", s)
	}

	n, err := strconv.Atoi(numStr)
	if err != nil || n < 0 {
		return Round{}, fmt.Errorf("invalid round number in %q", s)
	}

	return Round{Losers: losers, Number: n}, nil
}

// String returns the round as "WR2" or "LR1" format.
func (r Round) String() string {
	if r.Losers {
		return fmt.Sprintf("LR%d", r.Number)
	}
	return fmt.Sprintf("WR%d", r.Number)
}

// FromStartGG converts a start.gg round integer to a Round.
// start.gg uses positive numbers for winners and negative for losers.
func RoundFromStartGG(round int) Round {
	if round < 0 {
		return Round{Losers: true, Number: -round}
	}
	return Round{Losers: false, Number: round}
}
