package seed

import (
	"crypto/rand"
	"encoding/csv"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
)

type conflict struct {
	priority int
	players  []int // ID
}

const (
	CONFLICT_FILE                = "conflicts.csv"
	CONFLICT_RESOLUTION_ATTEMPTS = 1000
	CONFLICT_RESOLUTION_VARIANCE = 5
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

/* This is a Monte Carlo algorithm that randomly
 * swaps seeds in brackets and then evaluates the
 * score of that seeding. It keeps the randomly
 * generated seeding variant with the lowest
 * conflictScore.
 */
func resolveConflicts(bracket bracket, players []player) []player {
	conflicts := readConflictsFile()

	var best []player
	lowestScore := 999.0

	for range CONFLICT_RESOLUTION_ATTEMPTS {
		newPlayers := randomizeSeeds(players)

		conflictScore, _ := checkConflict(bracket, conflicts, newPlayers)

		seedDiffScore := 0.0
		for i, p := range newPlayers {
			for j, q := range players {
				if p.id == q.id {
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

		if conflictScore < lowestScore {
			lowestScore = conflictScore
			best = newPlayers
		}
	}

	log.Info("Seeds after conflict resolution")
	log.Printf("%-5s %-6s %25s %6s %-6s", "Seed", "Rating", "Name", "Change", "ID")
	log.Print("---------------------------------------------------------")
	for i, p := range best {
		for j, q := range players {
			if p == q {
				diff := j - i
				seed := i + 1
				if diff > 0 {
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.rating, p.name, "\033[32m↑", diff, "\033[0m", p.id)
				} else if diff < 0 {
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.rating, p.name, "\033[31m↓", -diff, "\033[0m", p.id)
				} else {
					log.Printf("%-5d %-6.1f %25s  %-6s %d", seed, p.rating, p.name, "", p.id)
				}
			}
		}
	}

	log.Debug("Finished conflict resolution", "score", lowestScore, "checks", CONFLICT_RESOLUTION_ATTEMPTS)

	return best
}

func readConflictsFile() []conflict {
	var conflicts []conflict
	file, err := os.Open(CONFLICT_FILE)
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
func randomizeSeeds(players []player) []player {
	newPlayers := make([]player, len(players))
	copy(newPlayers, players)
	for j := range players[3:] {
		i := j + 3
		r := rand.Reader
		x, _ := rand.Int(r, big.NewInt(CONFLICT_RESOLUTION_VARIANCE))
		if x.Int64() == 1 {
			temp := newPlayers[i]
			newPlayers[i] = newPlayers[i-1]
			newPlayers[i-1] = temp
		}
	}

	return newPlayers
}

/*
 * Returns the conflict sum, which matches the following
 * 4 * p1_conflicts + 3 * p2_conflicts + 2 * p3_conflicts
 */
func checkConflict(b bracket, cons []conflict, players []player) (float64, int) {
	conflictScore := 0.0
	conflictSum := 0

	for _, s := range b.sets {
		for _, con := range cons {
			if s == nil {
				break
			}
			p1 := players[s.player1-1].id
			p2 := players[s.player2-1].id
			if con.check(p1, p2) {
				conflictScore += float64(5 - con.priority)
				conflictSum += 1
			}
		}
	}

	return conflictScore, conflictSum
}
