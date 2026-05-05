package conflicts

import (
	"fmt"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

// printConflicts displays conflict information
func printConflicts(conflicts []Conflict) {
	for _, con := range conflicts {
		label := fmt.Sprintf("p%d", con.Priority)
		if con.isRequest() {
			label = fmt.Sprintf("r%d", -con.Priority)
		}
		roundStr := ""
		if con.Round != nil {
			roundStr = fmt.Sprintf(" %s", con.Round)
		}
		log.Printf("%-15s %15s  |  %s%s - %s", con.Players[0].Name, con.Players[1].Name, label, roundStr, con.Reason)
	}
}

// printUnresolved prints unresolved conflicts and unfulfilled set requests
// in a single pass, routing each to the formatter appropriate for its kind
// so nothing is listed twice.
func printUnresolved(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) {
	unresolved := listUnresolvedConflicts(bracket, conflicts, players)
	if len(unresolved) == 0 {
		return
	}

	var plainConflicts, setRequests []Conflict
	for _, c := range unresolved {
		if c.isRequest() {
			setRequests = append(setRequests, c)
		} else {
			plainConflicts = append(plainConflicts, c)
		}
	}

	printConflicts(plainConflicts)
	printSetRequests(bracket, setRequests, players)
}

// printSetRequests displays a set of unfulfilled set requests along with
// the round in which each requested pair is currently scheduled to meet
// (or "—" if they are not scheduled to meet at all). The caller is
// expected to pass only unresolved requests — this function doesn't
// re-check fulfillment.
func printSetRequests(bracket *brackets.Bracket, requests []Conflict, players []brackets.Player) {
	// Filter defensively so the header count matches what's actually
	// printed — callers may hand us a raw conflict list or slip a
	// malformed entry in.
	printable := make([]Conflict, 0, len(requests))
	for _, c := range requests {
		if c.isRequest() && len(c.Players) >= 2 {
			printable = append(printable, c)
		}
	}
	if len(printable) == 0 {
		return
	}

	roundMap := buildRoundMap(bracket)

	// Build an id→name index from the current player list so we can display
	// names even when the stored conflict doesn't have them populated.
	nameByID := make(map[any]string, len(players))
	for _, p := range players {
		nameByID[p.ID] = p.Name
	}

	nameFor := func(p Player) string {
		if p.Name != "" {
			return p.Name
		}
		if n, ok := nameByID[p.ID]; ok && n != "" {
			return n
		}
		return fmt.Sprintf("%v", p.ID)
	}

	log.Warnf("%d unfulfilled set request(s):", len(printable))
	for _, con := range printable {
		priorityLabel := fmt.Sprintf("r%d", -con.Priority)

		requestedRound := "any round"
		if con.Round != nil {
			requestedRound = con.Round.String()
		}

		// Find the round in which this pair currently meets (if any),
		// regardless of whether it matches the requested round — that
		// difference is the most useful diagnostic for the TO.
		currentRound := "—"
		for _, s := range bracket.Sets {
			if s == nil {
				continue
			}
			p1 := players[s.Player1-1].ID
			p2 := players[s.Player2-1].ID
			if !con.check(p1, p2) {
				continue
			}
			if r, ok := roundMap[s]; ok {
				currentRound = r.String()
			}
			break
		}

		log.Printf("  %-15s vs %-15s | %s | requested: %-9s | current: %s",
			nameFor(con.Players[0]), nameFor(con.Players[1]),
			priorityLabel, requestedRound, currentRound,
		)
	}
}

// printSeeds displays seed changes between original and final seeding
func printSeeds(before, after []brackets.Player) {
	log.Printf("%-5s %-6s %25s %6s %-7s", "Seed", "Rating", "Name", "Change", "ID")
	log.Print("---------------------------------------------------------")
	
	for i, p := range after {
		for j, q := range before {
			if p == q {
				diff := j - i
				seed := i + 1
				
				switch {
				case diff > 0:
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %s", seed, p.Rating, p.Name, "\033[32m↑", diff, "\033[0m", startgg.ToString(p.ID))
				case diff < 0:
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %s", seed, p.Rating, p.Name, "\033[31m↓", -diff, "\033[0m", startgg.ToString(p.ID))
				default:
					log.Printf("%-5d %-6.1f %25s  %-6s %s", seed, p.Rating, p.Name, "", startgg.ToString(p.ID))
				}
			}
		}
	}
}
