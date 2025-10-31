package conflicts

import (
	"math"
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/charmbracelet/log"
)

const (
	ConflictResolutionAttempts = 1000
	ConflictResolutionVariance = 6
)

type monteCarloResult struct {
	players []brackets.Player
	score   float64
}

// ResolveConflicts uses Monte Carlo algorithm to find optimal seeding
func ResolveConflicts(bracket *brackets.Bracket, conflicts []Conflict, players []brackets.Player) []brackets.Player {
	cache := newConflictCache(players, bracket)

	best := players
	lowestScore := calculateConflictScoreCached(cache, conflicts, players)

	log.Infof("conflictScore before resolution: %.2f", lowestScore)
	if unresolved := listUnresolvedConflicts(bracket, conflicts, players); len(unresolved) > 0 {
		printConflicts(unresolved)
	}

	if lowestScore > 0.0 {
		best, lowestScore = runParallelMonteCarlo(cache, conflicts, players, lowestScore)
	}

	logResults(players, best, bracket, conflicts, lowestScore)
	return best
}

// runParallelMonteCarlo executes the Monte Carlo algorithm with parallel workers
func runParallelMonteCarlo(cache *conflictCache, conflicts []Conflict, players []brackets.Player, initialScore float64) ([]brackets.Player, float64) {
	best := players
	lowestScore := initialScore
	numWorkers := runtime.NumCPU()

	for v := ConflictResolutionVariance; v > 2 && lowestScore > 0.0; v-- {
		attempts := calculateAttempts(v)
		log.Debug("Running parallel monte carlo", "variance", v, "attempts", attempts, "workers", numWorkers)

		result := runVarianceIteration(cache, conflicts, players, v, attempts, numWorkers, lowestScore)
		if result.score < lowestScore {
			log.Debug("Found new best", "variance", v, "score", result.score)
			lowestScore = result.score
			best = result.players

			if result.score == 0.0 {
				log.Info("Perfect solution found, terminating early")
				break
			}
		}

		// Adaptive variance reduction
		if lowestScore < 5.0 && v > 3 {
			v-- // Skip a variance level for good solutions
		}
	}

	return best, lowestScore
}

// runVarianceIteration runs one iteration of the Monte Carlo algorithm
func runVarianceIteration(cache *conflictCache, conflicts []Conflict, players []brackets.Player, variance, attempts, numWorkers int, currentBest float64) monteCarloResult {
	resultChan := make(chan monteCarloResult, numWorkers)
	attemptsPerWorker := attempts / numWorkers

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerBest := players
			workerLowest := currentBest

			for range attemptsPerWorker {
				newPlayers := randomizeSeeds(players, variance)
				score := calculateConflictScoreCached(cache, conflicts, newPlayers)

				if score < workerLowest {
					workerLowest = score
					workerBest = make([]brackets.Player, len(newPlayers))
					copy(workerBest, newPlayers)
				}
			}

			resultChan <- monteCarloResult{workerBest, workerLowest}
		}()
	}

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	bestResult := monteCarloResult{players, currentBest}
	for result := range resultChan {
		if result.score < bestResult.score {
			bestResult = result
		}
	}

	return bestResult
}

// randomizeSeeds randomly shifts seeding (doesn't impact seeds 1 and 2)
func randomizeSeeds(players []brackets.Player, variance int) []brackets.Player {
	newPlayers := make([]brackets.Player, len(players))
	copy(newPlayers, players)

	for j := range players[3:] {
		i := j + 3
		if rand.IntN(variance) == 0 {
			temp := newPlayers[i]
			n := 1
			newPlayers[i] = newPlayers[i-n]
			newPlayers[i-n] = temp
		}
	}

	return newPlayers
}

// calculateAttempts determines number of attempts for given variance
func calculateAttempts(variance int) int {
	return int(math.Max(
		ConflictResolutionAttempts*math.Pow(float64(6-variance), 1.5),
		ConflictResolutionAttempts,
	))
}

// logResults prints the final results
func logResults(original, final []brackets.Player, bracket *brackets.Bracket, conflicts []Conflict, score float64) {
	log.Info("Seeds after conflict resolution", "score", score)
	printSeeds(original, final)

	_, numConflicts := checkConflict(bracket, conflicts, final)
	if numConflicts > 0 {
		log.Warnf("%d conflicts were unresolved", numConflicts)
		if unresolved := listUnresolvedConflicts(bracket, conflicts, final); len(unresolved) > 0 {
			printConflicts(unresolved)
		}
	}

	log.Debug("Finished conflict resolution", "score", score)
}
