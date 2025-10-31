// Package validation provides tournament set score validation for best-of-3 and best-of-5 formats.
package validation

import (
	"fmt"

	"github.com/AaronLieb/octagon/tournament"
)

func ValidateSetScore(gameResults []tournament.GameResult) error {
	p1Wins := 0
	p2Wins := 0
	for _, result := range gameResults {
		switch result.Winner {
		case 1:
			p1Wins++
		case 2:
			p2Wins++
		}
	}

	// Validate Bo3 or Bo5 scores
	if (p1Wins != 2 || p2Wins > 1) && (p2Wins != 2 || p1Wins > 1) &&
		(p1Wins != 3 || p2Wins > 2) && (p2Wins != 3 || p1Wins > 2) {
		return fmt.Errorf("invalid score: must be best-of-3 (2-0, 2-1) or best-of-5 (3-0, 3-1, 3-2)")
	}

	return nil
}
