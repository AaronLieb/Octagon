package conflicts

import (
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/charmbracelet/log"
)

const (
	ConflictResolutionAttempts = 5000
	ConflictResolutionVariance = 8
	AttemptsAddedPerRound      = 5000
	MinVariance                = 3
	MaxAttemptsPerRound        = 30000
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

	for v := ConflictResolutionVariance; v > MinVariance && lowestScore > 0.0; v-- {
		attempts := calculateAttempts(v)
		log.Debug("Running parallel monte carlo", "variance", v, "attempts", attempts, "workers", numWorkers)

		// Pass the best-so-far as the starting point for this iteration, so
		// progress compounds across variance levels rather than restarting
		// from the original seeding every time.
		result := runVarianceIteration(cache, conflicts, players, best, v, attempts, numWorkers, lowestScore)
		if result.score < lowestScore {
			log.Debug("Found new best", "variance", v, "score", result.score)
			lowestScore = result.score
			best = result.players

			if result.score == 0.0 {
				log.Info("Perfect solution found, terminating early")
				break
			}
		}
	}

	return best, lowestScore
}

// restartInterval controls how often a worker resets its walk back to the
// original seeding. Restarts inject diversity and help the MC escape local
// minima when perturbing only from the current best would stagnate.
const restartInterval = 200

// runVarianceIteration runs one iteration of the Monte Carlo algorithm.
//
// Within each worker we hill-climb: successive perturbations are applied to
// the current best rather than always starting from the original seeding.
// This lets progress compound and is essential for scenarios where the
// optimal solution requires several coordinated seed moves (e.g. multiple
// simultaneous set requests). We periodically restart from the original to
// escape local minima.
func runVarianceIteration(cache *conflictCache, conflicts []Conflict, original, startBest []brackets.Player, variance, attempts, numWorkers int, currentBest float64) monteCarloResult {
	resultChan := make(chan monteCarloResult, numWorkers)
	attemptsPerWorker := attempts / numWorkers

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			workerBest := make([]brackets.Player, len(startBest))
			copy(workerBest, startBest)
			workerLowest := currentBest

			walkFrom := workerBest
			sinceImprovement := 0

			for i := range attemptsPerWorker {
				// Periodic restart: alternate between walking from the current
				// best (exploitation) and resetting to the original seeding
				// (exploration). This prevents a worker from getting permanently
				// stuck in a local minimum.
				if i > 0 && sinceImprovement >= restartInterval {
					walkFrom = original
					sinceImprovement = 0
				}

				newPlayers := randomizeSeeds(walkFrom, variance)
				score := calculateConflictScoreCached(cache, conflicts, newPlayers)

				if score < workerLowest {
					workerLowest = score
					workerBest = make([]brackets.Player, len(newPlayers))
					copy(workerBest, newPlayers)
					walkFrom = workerBest
					sinceImprovement = 0
				} else {
					sinceImprovement++
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

	bestResult := monteCarloResult{startBest, currentBest}
	for result := range resultChan {
		if result.score < bestResult.score {
			bestResult = result
		}
	}

	return bestResult
}

// randomizeSeeds perturbs the seeding while preserving seeds 1 and 2.
// It combines two kinds of perturbation:
//
//  1. Per-position adjacent swaps (the original mechanism): each position
//     from index 3 onwards has a 1/variance chance of swapping with its
//     immediate predecessor. This produces mostly-local perturbations that
//     fine-tune the seeding once it's close to optimal.
//  2. An occasional long-range swap between two random non-reserved
//     positions. Without this, candidates that require moving a player
//     more than a couple of seeds are effectively unreachable in reasonable
//     attempt counts — which was why set requests on larger brackets
//     silently failed to resolve.
func randomizeSeeds(players []brackets.Player, variance int) []brackets.Player {
	newPlayers := make([]brackets.Player, len(players))
	copy(newPlayers, players)

	// Seeds 1 and 2 are preserved; we shuffle from index 2 onwards.
	// Brackets smaller than 4 have nothing meaningful to shuffle.
	if len(players) < 4 {
		return newPlayers
	}

	// Long-range swap. Gated by variance so lower variance (more aggressive
	// MC iteration) produces big jumps more often. We draw two indices in
	// [2, len-1] and swap them unconditionally if they differ.
	if rand.IntN(variance) == 0 {
		shufStart := 2
		n := len(players) - shufStart
		i := shufStart + rand.IntN(n)
		j := shufStart + rand.IntN(n)
		if i != j {
			newPlayers[i], newPlayers[j] = newPlayers[j], newPlayers[i]
		}
	}

	// Adjacent swaps (the original perturbation).
	for k := range players[3:] {
		i := k + 3
		if rand.IntN(variance) == 0 {
			newPlayers[i], newPlayers[i-1] = newPlayers[i-1], newPlayers[i]
		}
	}

	return newPlayers
}

// calculateAttempts determines number of attempts for given variance
func calculateAttempts(variance int) int {
	extraRounds := ConflictResolutionVariance - variance
	return min(int(ConflictResolutionAttempts+AttemptsAddedPerRound*extraRounds), MaxAttemptsPerRound)
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
