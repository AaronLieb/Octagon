package conflicts

import (
	"encoding/json"
	"math"
	"math/rand/v2"
	"os"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

const (
	CONFLICT_FILE                = "conflicts.json"
	CONFLICT_RESOLUTION_ATTEMPTS = 1000
	CONFLICT_RESOLUTION_VARIANCE = 6
)

// returns true if p1 and p2 are in the conflict
func (con *conflict) check(p1 startgg.ID, p2 startgg.ID) bool {
	flag := true
	for _, p := range con.Players {
		if startgg.ToString(p.Id) == startgg.ToString(p1) || startgg.ToString(p.Id) == startgg.ToString(p2) {
			if flag {
				flag = false
			} else {
				return true
			}
		}
	}
	return false
}

func printConflicts(cons []conflict) {
	for _, con := range cons {
		log.Printf("%-15s %15s\np%d - %s\n", con.Players[0].Name, con.Players[1].Name, con.Priority, con.Reason)
	}
}

func calculateConflictScore(bracket *brackets.Bracket, conflicts []conflict, players []brackets.Player, newPlayers []brackets.Player) float64 {
	conflictScore, _ := checkConflict(bracket, conflicts, newPlayers)

	// TODO: Prune this once the conflict score gets above certain value
	seedDiffScore := 0.0
	for i, p := range newPlayers {
		for j, q := range players {
			if p.Id == q.Id {
				// Low importance means changing the seed has less significance
				importance := 32.0 / math.Pow(math.Log2(float64(j)), 3)
				// keep inbetween [0.25, 1]
				importance = math.Min(1, math.Max(0.25, importance))
				seedDiff := math.Abs(float64(i-j)) * importance
				seedDiffScore += math.Pow(seedDiff, 1.5)
			}
		}
	}
	conflictScore += seedDiffScore / 2
	return conflictScore
}

/* This is a Monte Carlo algorithm that randomly
 * swaps seeds in brackets and then evaluates the
 * score of that seeding. It keeps the randomly
 * generated seeding variant with the lowest
 * conflictScore.
 */
func ResolveConflicts(bracket *brackets.Bracket, conflicts []conflict, players []brackets.Player) []brackets.Player {
	best := players
	lowestScore := calculateConflictScore(bracket, conflicts, players, players)

	log.Infof("conflictScore before resolution: %.2f", lowestScore)
	unresolved := listUnresolvedConflicts(bracket, conflicts, players)
	printConflicts(unresolved)

	if lowestScore != 0.0 {
		for v := CONFLICT_RESOLUTION_VARIANCE; v > 2; v-- {
			attempts := int(math.Max(CONFLICT_RESOLUTION_ATTEMPTS*math.Pow(float64(6-v), 1.5), CONFLICT_RESOLUTION_ATTEMPTS))
			log.Debug("Running monte carlo algo", "variance", v, "attempts", attempts)
			for range attempts {
				newPlayers := randomizeSeeds(players, v)

				conflictScore := calculateConflictScore(bracket, conflicts, players, newPlayers)

				if conflictScore < lowestScore {
					log.Debug("Found new best", "variance", v)
					lowestScore = conflictScore
					best = newPlayers
				}
			}
		}
	}

	log.Info("Seeds after conflict resolution", "score", lowestScore)
	printSeeds(players, best)
	_, numConflicts := checkConflict(bracket, conflicts, best)
	if numConflicts != 0 {
		log.Warnf("%d conflicts were unresolved", numConflicts)
		unresolved := listUnresolvedConflicts(bracket, conflicts, best)
		printConflicts(unresolved)
	}
	log.Debug("Finished conflict resolution", "score", lowestScore, "checks", CONFLICT_RESOLUTION_ATTEMPTS)

	return best
}

func printSeeds(before []brackets.Player, after []brackets.Player) {
	log.Printf("%-5s %-6s %25s %6s %-7s", "Seed", "Rating", "Name", "Change", "ID")
	log.Print("---------------------------------------------------------")
	for i, p := range after {
		for j, q := range before {
			if p == q {
				diff := j - i
				seed := i + 1
				if diff > 0 {
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.Rating, p.Name, "\033[32m↑", diff, "\033[0m", p.Id)
				} else if diff < 0 {
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.Rating, p.Name, "\033[31m↓", -diff, "\033[0m", p.Id)
				} else {
					log.Printf("%-5d %-6.1f %25s  %-6s %d", seed, p.Rating, p.Name, "", p.Id)
				}
			}
		}
	}
}

func GetConflicts(conflictFiles []string) []conflict {
	conflicts := readConflictsFile(CONFLICT_FILE)
	for _, file := range conflictFiles {
		conflicts = append(conflicts, readConflictsFile(file)...)
	}
	return conflicts
}

func readConflictsFile(fileName string) []conflict {
	var cons []conflict
	path := config.GetConfigPath() + fileName
	file, err := os.ReadFile(path)
	if err != nil {
		log.Warn("unable to read conflicts file: %v", err)
		return cons
	}

	err = json.Unmarshal(file, &cons)
	if err != nil {
		log.Errorf("unable to unmarshal conflicts file: %v", err)
		return cons
	}

	log.Info("Reading conflicts", "path", path, "n", len(cons))
	if len(cons) == 0 {
		log.Warn("No conflicts found in conflict file")
	}
	return cons
}

/*
* Randomly shifts seeding
* Doesn't impact seeds 1 and 2
 */
func randomizeSeeds(players []brackets.Player, variance int) []brackets.Player {
	newPlayers := make([]brackets.Player, len(players))
	copy(newPlayers, players)
	for j := range players[3:] {
		i := j + 3
		if rand.IntN(variance) == 0 {
			temp := newPlayers[i]
			newPlayers[i] = newPlayers[i-1]
			newPlayers[i-1] = temp
		}
	}

	return newPlayers
}

/*
 * Returns the conflict sum, lower value is better.
 * If a conflict is unresolved, it will add to the sum.
 * The higher the priority of the conflict the more it
 * adds to the sum
 */
func checkConflict(b *brackets.Bracket, cons []conflict, players []brackets.Player) (float64, int) {
	conflictScore := 0.0
	conflictSum := 0

	for _, s := range b.Sets {
		for _, con := range cons {
			if s == nil {
				break
			}
			p1 := players[s.Player1-1].Id
			p2 := players[s.Player2-1].Id
			if con.check(p1, p2) {
				conflictScore += float64(2 + con.Priority)
				conflictSum += 1
			}
		}
	}

	return conflictScore, conflictSum
}

/*
 * Returns a list of conflicts that were unresolved
 */
func listUnresolvedConflicts(b *brackets.Bracket, cons []conflict, players []brackets.Player) []conflict {
	var unresolved []conflict
	for _, s := range b.Sets {
		for _, con := range cons {
			if s == nil {
				break
			}
			p1 := players[s.Player1-1].Id
			p2 := players[s.Player2-1].Id
			if con.check(p1, p2) {
				unresolved = append(unresolved, con)
			}
		}
	}
	return unresolved
}
