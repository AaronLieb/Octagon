package conflicts

import (
	"encoding/csv"
	"math"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/AaronLieb/octagon/bracket"
	"github.com/charmbracelet/log"
)

const (
	CONFLICT_FILE                = "conflicts.csv"
	CONFLICT_RESOLUTION_ATTEMPTS = 0
	CONFLICT_RESOLUTION_VARIANCE = 3
)

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

func calculateConflictScore(bracket bracket.Bracket, conflicts []conflict, players []bracket.Player, newPlayers []bracket.Player) float64 {
	conflictScore, _ := checkConflict(bracket, conflicts, newPlayers)

	seedDiffScore := 0.0
	for i, p := range newPlayers {
		for j, q := range players {
			if p.Id == q.Id {
				// Low importance means changing the seed has less significance
				importance := 32.0 / math.Pow(math.Log2(float64(j)), 3)
				// keep inbetween [0.25, 1]
				importance = math.Max(1, math.Min(0.25, importance))
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
func ResolveConflicts(bracket bracket.Bracket, conflicts []conflict, players []bracket.Player) []bracket.Player {
	best := players
	lowestScore := calculateConflictScore(bracket, conflicts, players, players)

	log.Infof("conflictScore before resolution: %.2f", lowestScore)

	if lowestScore != 0.0 {
		for range CONFLICT_RESOLUTION_ATTEMPTS {
			newPlayers := randomizeSeeds(players)

			conflictScore := calculateConflictScore(bracket, conflicts, players, newPlayers)

			if conflictScore < lowestScore {
				lowestScore = conflictScore
				best = newPlayers
			}
		}
	}

	log.Info("Seeds after conflict resolution", "score", lowestScore)
	printSeeds(players, best)
	_, numConflicts := checkConflict(bracket, conflicts, best)
	// if numConflicts != 0 {
	log.Warnf("%d conflicts were unresolved", numConflicts)
	// }
	log.Debug("Finished conflict resolution", "score", lowestScore, "checks", CONFLICT_RESOLUTION_ATTEMPTS)

	return best
}

func printSeeds(before []bracket.Player, after []bracket.Player) {
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
	var conflicts []conflict
	for _, file := range conflictFiles {
		conflicts = append(conflicts, readConflictsFile(file)...)
	}
	return conflicts
}

func readConflictsFile(fileName string) []conflict {
	var conflicts []conflict
	file, err := os.Open(fileName)
	if err != nil {
		log.Warn("unable to find or open conflicts", "file", CONFLICT_FILE)
		return conflicts
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Reading conflicts", "file", CONFLICT_FILE, "n", len(records))
	conflicts = make([]conflict, len(records))
	for i, row := range records {
		priority, err := strconv.Atoi(row[0])
		if err != nil {
			log.Fatalf("unable to parse priority in line %d of conflicts file: %v\n", i, err)
		}
		players := make([]int, len(row[1:]))
		for j, s := range row[1:] {
			p, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("unable to parse player id in line %d of conflicts file: %v\n", i, err)
			}
			players[j] = p
		}
		conflicts[i] = conflict{
			priority: priority,
			players:  players,
		}
	}
	if len(conflicts) == 0 {
		log.Warn("No conflicts found in conflict file")
	}
	return conflicts
}

/*
* Randomly shifts seeding
* Doesn't impact seeds 1 and 2
 */
func randomizeSeeds(players []bracket.Player) []bracket.Player {
	newPlayers := make([]bracket.Player, len(players))
	copy(newPlayers, players)
	for j := range players[3:] {
		i := j + 3
		if rand.IntN(CONFLICT_RESOLUTION_VARIANCE) == 0 {
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
func checkConflict(b bracket.Bracket, cons []conflict, players []bracket.Player) (float64, int) {
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
				conflictScore += float64(2 + con.priority)
				conflictSum += 1
			}
		}
	}

	return conflictScore, conflictSum
}
