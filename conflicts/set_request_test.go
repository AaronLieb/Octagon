package conflicts

import (
	"fmt"
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

// makePlayers returns n players with IDs 1..n and descending ratings.
func makePlayers(n int) []brackets.Player {
	players := make([]brackets.Player, n)
	for i := range players {
		players[i] = brackets.Player{
			ID:     int64(i + 1),
			Name:   fmt.Sprintf("P%d", i+1),
			Rating: float64(n - i),
		}
	}
	return players
}

// seedOf returns the 1-indexed seed of the player with the given ID in players.
func seedOf(players []brackets.Player, id int64) int {
	for i, p := range players {
		if p.ID == id {
			return i + 1
		}
	}
	return -1
}

// playersMeet returns the set where p1 and p2 meet in the bracket, along with
// the round they meet in. It returns (nil, Round{}, false) if they don't meet.
func playersMeet(bracket *brackets.Bracket, players []brackets.Player, id1, id2 int64) (*brackets.Set, brackets.Round, bool) {
	s1 := seedOf(players, id1)
	s2 := seedOf(players, id2)
	if s1 < 0 || s2 < 0 {
		return nil, brackets.Round{}, false
	}
	roundMap := buildRoundMap(bracket)
	for _, s := range bracket.Sets {
		if s == nil {
			continue
		}
		if (s.Player1 == s1 && s.Player2 == s2) || (s.Player1 == s2 && s.Player2 == s1) {
			return s, roundMap[s], true
		}
	}
	return nil, brackets.Round{}, false
}

// TestCachedScorerMatchesReference asserts that the cached scorer produces the
// same (score, count) as the reference scorer for a variety of inputs. This is
// the invariant that makes the Monte Carlo algorithm work: the cached score is
// what the optimizer actually minimizes, and if it disagrees with the reference
// the optimizer has no reliable signal.
func TestCachedScorerMatchesReference(t *testing.T) {
	type scenario struct {
		name      string
		size      int
		conflicts func(players []brackets.Player) []Conflict
	}

	scenarios := []scenario{
		{
			name: "no conflicts",
			size: 8,
			conflicts: func(_ []brackets.Player) []Conflict {
				return nil
			},
		},
		{
			name: "single conflict triggered",
			size: 8,
			conflicts: func(_ []brackets.Player) []Conflict {
				return []Conflict{{
					Priority: 2,
					Players:  []Player{{ID: int64(1)}, {ID: int64(8)}},
				}}
			},
		},
		{
			name: "single set request unfulfilled",
			size: 8,
			conflicts: func(_ []brackets.Player) []Conflict {
				// P1 vs P6 don't meet in default 8-player seeding
				return []Conflict{{
					Priority: -2,
					Players:  []Player{{ID: int64(1)}, {ID: int64(6)}},
				}}
			},
		},
		{
			name: "single set request fulfilled",
			size: 8,
			conflicts: func(_ []brackets.Player) []Conflict {
				// P1 vs P8 meet in WR1
				return []Conflict{{
					Priority: -2,
					Players:  []Player{{ID: int64(1)}, {ID: int64(8)}},
				}}
			},
		},
		{
			name: "multiple unfulfilled requests with different priorities",
			size: 16,
			conflicts: func(_ []brackets.Player) []Conflict {
				return []Conflict{
					{Priority: -1, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
					{Priority: -2, Players: []Player{{ID: int64(1)}, {ID: int64(7)}}},
					{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(10)}}},
				}
			},
		},
		{
			name: "round-specific request unfulfilled",
			size: 16,
			conflicts: func(_ []brackets.Player) []Conflict {
				// P1 vs P2 only meet in the finals, not WR1
				return []Conflict{{
					Priority: -2,
					Round:    wr(1),
					Players:  []Player{{ID: int64(1)}, {ID: int64(2)}},
				}}
			},
		},
		{
			name: "round-specific request fulfilled",
			size: 16,
			conflicts: func(_ []brackets.Player) []Conflict {
				// P1 vs P16 meet in WR1
				return []Conflict{{
					Priority: -2,
					Round:    wr(1),
					Players:  []Player{{ID: int64(1)}, {ID: int64(16)}},
				}}
			},
		},
		{
			name: "mixed conflicts and requests",
			size: 16,
			conflicts: func(_ []brackets.Player) []Conflict {
				return []Conflict{
					{Priority: 3, Players: []Player{{ID: int64(1)}, {ID: int64(16)}}},
					{Priority: -1, Players: []Player{{ID: int64(2)}, {ID: int64(5)}}},
					{Priority: 1, Round: wr(1), Players: []Player{{ID: int64(8)}, {ID: int64(9)}}},
				}
			},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			players := makePlayers(sc.size)
			bracket := brackets.CreateBracket(sc.size)
			cons := sc.conflicts(players)

			refScore, refCount := checkConflict(bracket, cons, players)
			cache := newConflictCache(players, bracket)
			cachedScore, cachedCount := checkConflictCached(cache, cons, players)

			if refScore != cachedScore || refCount != cachedCount {
				t.Errorf("cached scorer mismatch: reference=(%.1f, %d) cached=(%.1f, %d)",
					refScore, refCount, cachedScore, cachedCount)
			}
		})
	}
}

// TestNoPenaltyWithoutRequests is a regression test for a bug where the cached
// scorer unconditionally added a penalty of 5 at the end, even when there were
// zero set requests in the conflict list.
func TestNoPenaltyWithoutRequests(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	cache := newConflictCache(players, bracket)
	score, count := checkConflictCached(cache, nil, players)

	if score != 0.0 || count != 0 {
		t.Errorf("expected (0.0, 0) with no conflicts, got (%.1f, %d)", score, count)
	}

	// Also with a conflict list that contains a regular (non-request) conflict
	// that isn't triggered.
	cons := []Conflict{{
		Priority: 2,
		Players:  []Player{{ID: int64(1)}, {ID: int64(4)}}, // They don't meet in WR1
		Round:    wr(1),
	}}
	score, count = checkConflictCached(cache, cons, players)
	if score != 0.0 || count != 0 {
		t.Errorf("expected (0.0, 0) with no triggered conflicts, got (%.1f, %d)", score, count)
	}
}

// TestEachRequestCountedSeparately is a regression test for a bug where the
// cached scorer collapsed all unfulfilled requests into a single penalty of 5.
func TestEachRequestCountedSeparately(t *testing.T) {
	players := makePlayers(16)
	bracket := brackets.CreateBracket(16)

	cons := []Conflict{
		{Priority: -1, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
		{Priority: -2, Players: []Player{{ID: int64(1)}, {ID: int64(7)}}},
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(10)}}},
	}

	cache := newConflictCache(players, bracket)
	score, count := checkConflictCached(cache, cons, players)

	// None of these pairs meet in the default 16-player seeding, so all 3 should
	// penalize. Penalties: (2+1) + (2+2) + (2+3) = 12. Count: 3.
	if count != 3 {
		t.Errorf("expected 3 unfulfilled requests, got %d", count)
	}
	if score != 12.0 {
		t.Errorf("expected score 12.0, got %.1f", score)
	}
}

// TestPenaltyScalesWithPriority is a regression test for a bug where the cached
// scorer hardcoded the penalty at 5 regardless of the request's priority.
func TestPenaltyScalesWithPriority(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)
	cache := newConflictCache(players, bracket)

	for _, priority := range []int{-1, -2, -3} {
		cons := []Conflict{{
			Priority: priority,
			Players:  []Player{{ID: int64(1)}, {ID: int64(6)}}, // P1 vs P6 don't meet
		}}
		score, _ := checkConflictCached(cache, cons, players)
		expected := 2.0 + float64(-priority)
		if score != expected {
			t.Errorf("priority=%d: expected score %.1f, got %.1f", priority, expected, score)
		}
	}
}

// TestRoundIndexingMatchesParseRound is a regression test for the 0/1 indexing
// bug. A user typing --round WR1 should match the first winners round.
func TestRoundIndexingMatchesParseRound(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	round, err := brackets.ParseRound("WR1")
	if err != nil {
		t.Fatalf("ParseRound(WR1): %v", err)
	}

	// P1 vs P8 meet in the first winners round. A request with that round
	// (as the CLI would parse it) should be fulfilled by the default seeding.
	cons := []Conflict{{
		Priority: -1,
		Round:    &round,
		Players:  []Player{{ID: int64(1)}, {ID: int64(8)}},
	}}

	score, count := checkConflict(bracket, cons, players)
	if score != 0.0 || count != 0 {
		t.Errorf("WR1 request for P1 vs P8 should be fulfilled, got score=%.1f count=%d", score, count)
	}

	cache := newConflictCache(players, bracket)
	cachedScore, cachedCount := checkConflictCached(cache, cons, players)
	if cachedScore != 0.0 || cachedCount != 0 {
		t.Errorf("cached: WR1 request for P1 vs P8 should be fulfilled, got score=%.1f count=%d", cachedScore, cachedCount)
	}
}

// TestRoundIndexingLosersRound verifies the 1-indexed round numbering also
// matches losers-side rounds parsed from the CLI.
func TestRoundIndexingLosersRound(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// In an 8-player bracket, LR1 has 2 sets. Find a pair that meets there.
	if len(bracket.LosersRounds) == 0 || len(bracket.LosersRounds[0]) == 0 {
		t.Skip("no losers round 1 sets in this bracket")
	}
	firstLR1Set := bracket.LosersRounds[0][0]
	id1 := players[firstLR1Set.Player1-1].ID
	id2 := players[firstLR1Set.Player2-1].ID

	round, err := brackets.ParseRound("LR1")
	if err != nil {
		t.Fatalf("ParseRound(LR1): %v", err)
	}

	cons := []Conflict{{
		Priority: -1,
		Round:    &round,
		Players:  []Player{{ID: id1}, {ID: id2}},
	}}

	score, count := checkConflict(bracket, cons, players)
	if score != 0.0 || count != 0 {
		t.Errorf("LR1 request for a pair that meets in LR1 should be fulfilled, got score=%.1f count=%d", score, count)
	}
}

// resolvesRequest is a helper that runs ResolveConflicts and asserts that
// every negative-priority request is fulfilled in the returned seeding.
func resolvesRequest(t *testing.T, size int, requests []Conflict) []brackets.Player {
	t.Helper()
	players := makePlayers(size)
	bracket := brackets.CreateBracket(size)

	result := ResolveConflicts(bracket, requests, players)

	unresolved := listUnresolvedConflicts(bracket, requests, result)
	for _, u := range unresolved {
		if u.isRequest() {
			names := []string{}
			for _, p := range u.Players {
				names = append(names, fmt.Sprintf("ID=%v", p.ID))
			}
			t.Errorf("set request %v (priority=%d, round=%v) unfulfilled after resolution", names, u.Priority, u.Round)
		}
	}

	return result
}

func TestResolveSetRequest_8Players(t *testing.T) {
	// P1 vs P6 don't meet in default 8-player seeding (WR1: 1v8, 4v5, 2v7, 3v6;
	// WR2: 1v4, 2v3; WR3: 1v2). MC must rearrange seeds so they meet somewhere.
	resolvesRequest(t, 8, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
	})
}

func TestResolveSetRequest_16Players(t *testing.T) {
	// P1 vs P10 don't meet in default 16-player seeding.
	resolvesRequest(t, 16, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(10)}}},
	})
}

func TestResolveSetRequest_32Players(t *testing.T) {
	// P1 vs P20 don't meet in default 32-player seeding.
	resolvesRequest(t, 32, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(20)}}},
	})
}

func TestResolveMultipleSetRequests(t *testing.T) {
	// Three independent requests at 16 players. All pairs involve distinct
	// low seeds that don't meet in the default seeding.
	resolvesRequest(t, 16, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
		{Priority: -3, Players: []Player{{ID: int64(3)}, {ID: int64(10)}}},
		{Priority: -3, Players: []Player{{ID: int64(4)}, {ID: int64(11)}}},
	})
}

func TestResolveSetRequestWithRoundWR1(t *testing.T) {
	// Request P1 vs P6 specifically in WR1 (16-player bracket). WR1 pairs are
	// high-vs-low seed (1v16, 8v9, etc.), so for P1 vs P6 to play WR1 the MC
	// must move P6 to seed 16 (or P1's opponent's seed).
	players := makePlayers(16)
	bracket := brackets.CreateBracket(16)

	requests := []Conflict{
		{Priority: -3, Round: wr(1), Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
	}

	result := ResolveConflicts(bracket, requests, players)

	// Verify P1 and P6 meet in WR1 specifically.
	_, round, meet := playersMeet(bracket, result, int64(1), int64(6))
	if !meet {
		t.Fatal("P1 and P6 don't meet at all after resolution")
	}
	if round.Losers || round.Number != 1 {
		t.Errorf("P1 and P6 should meet in WR1, got %s", round)
	}
}

func TestResolveSetRequestAlreadyFulfilled(t *testing.T) {
	// P1 vs P8 meet in WR1 in default 8-player seeding. A request for them
	// should be a no-op, and the MC should leave seeding essentially unchanged.
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	requests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
	}

	result := ResolveConflicts(bracket, requests, players)

	// Request must still be fulfilled.
	_, count := checkConflict(bracket, requests, result)
	if count != 0 {
		t.Errorf("already-fulfilled request should remain fulfilled, got count=%d", count)
	}

	// Top 2 seeds should be unchanged (they're reserved anyway).
	if result[0].ID != int64(1) || result[1].ID != int64(2) {
		t.Errorf("top seeds should be preserved, got %v %v", result[0].ID, result[1].ID)
	}
}

func TestResolveMixedConflictsAndRequests(t *testing.T) {
	// A conflict (P1 vs P8 shouldn't play WR1) and a request (P1 vs P6 should
	// play somewhere). The MC must handle both simultaneously.
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	cons := []Conflict{
		// Conflict: avoid P1 vs P8 in WR1 (they currently meet there).
		{Priority: 3, Round: wr(1), Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
		// Request: P1 vs P6 should meet somewhere.
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
	}

	result := ResolveConflicts(bracket, cons, players)

	unresolved := listUnresolvedConflicts(bracket, cons, result)
	if len(unresolved) != 0 {
		for _, u := range unresolved {
			t.Errorf("unresolved: priority=%d round=%v players=%v", u.Priority, u.Round, u.Players)
		}
	}
}

func TestResolveConflictsTwoPlayersNoCrash(t *testing.T) {
	// Regression test: randomizeSeeds used to panic with fewer than 3 players.
	players := makePlayers(2)
	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{{Player1: 1, Player2: 2}},
		WinnersRounds: [][]*brackets.Set{
			{{Player1: 1, Player2: 2}},
		},
	}

	// Attach the same set pointer so round map works.
	bracket.Sets = bracket.WinnersRounds[0]

	// Run with a conflict that forces Monte Carlo to execute.
	cons := []Conflict{{
		Priority: 3,
		Players:  []Player{{ID: int64(1)}, {ID: int64(2)}},
	}}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ResolveConflicts panicked with 2 players: %v", r)
		}
	}()

	result := ResolveConflicts(bracket, cons, players)
	if len(result) != 2 {
		t.Errorf("expected 2 players back, got %d", len(result))
	}
}

func TestResolveConflictsThreePlayersNoCrash(t *testing.T) {
	// Regression test: randomizeSeeds used to panic with fewer than 3 players.
	// 3 players should also not panic even though there's nothing to shuffle.
	players := makePlayers(3)
	bracket := brackets.CreateBracket(3)

	// Use a conflict that will trigger MC.
	cons := []Conflict{{
		Priority: 3,
		Players:  []Player{{ID: int64(1)}, {ID: int64(3)}},
	}}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ResolveConflicts panicked with 3 players: %v", r)
		}
	}()

	_ = ResolveConflicts(bracket, cons, players)
}

func TestRandomizeSeedsTinyBracket(t *testing.T) {
	// Unit-level regression for randomizeSeeds itself.
	for _, n := range []int{0, 1, 2, 3} {
		players := makePlayers(n)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("randomizeSeeds panicked with %d players: %v", n, r)
				}
			}()
			got := randomizeSeeds(players, 2)
			if len(got) != n {
				t.Errorf("randomizeSeeds with %d players returned %d", n, len(got))
			}
		}()
	}
}

// TestResolvePreservesTopSeeds verifies the reserved top-2-seeds invariant
// holds after Monte Carlo for a range of bracket sizes and request patterns.
func TestResolvePreservesTopSeeds(t *testing.T) {
	sizes := []int{8, 16, 32}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("%d-players", size), func(t *testing.T) {
			players := makePlayers(size)
			bracket := brackets.CreateBracket(size)

			// Request that the lowest seed meets the highest seed (rearrange-heavy).
			requests := []Conflict{{
				Priority: -3,
				Players:  []Player{{ID: int64(3)}, {ID: int64(size - 2)}},
			}}

			result := ResolveConflicts(bracket, requests, players)

			if result[0].ID != int64(1) {
				t.Errorf("seed 1 changed: got ID=%v", result[0].ID)
			}
			if result[1].ID != int64(2) {
				t.Errorf("seed 2 changed: got ID=%v", result[1].ID)
			}
		})
	}
}

func TestResolveSetRequest_64Players(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 64-player MC stress test in short mode")
	}
	// 64-player bracket, request a pair that requires significant rearrangement.
	resolvesRequest(t, 64, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(30)}}},
	})
}

func TestResolveSetRequestLosersRoundSpecific(t *testing.T) {
	// Request a pair meet in LR2 specifically (16-player bracket).
	players := makePlayers(16)
	bracket := brackets.CreateBracket(16)

	// Pick a pair that doesn't meet in LR2 by default. We'll let the MC
	// move players to make them meet there.
	requests := []Conflict{
		{Priority: -3, Round: lr(2), Players: []Player{{ID: int64(3)}, {ID: int64(12)}}},
	}

	result := ResolveConflicts(bracket, requests, players)

	unresolved := listUnresolvedConflicts(bracket, requests, result)
	if len(unresolved) > 0 {
		t.Errorf("LR2 request unfulfilled after resolution")
		_, round, meet := playersMeet(bracket, result, int64(3), int64(12))
		if meet {
			t.Logf("players do meet but in %s instead of LR2", round)
		}
	}
}

// TestResolveSetRequestStability runs the resolver multiple times to catch
// flakiness from randomized Monte Carlo. A working implementation should
// fulfill a basic request reliably on every run.
func TestResolveSetRequestStability(t *testing.T) {
	const runs = 10
	for i := 0; i < runs; i++ {
		players := makePlayers(16)
		bracket := brackets.CreateBracket(16)
		requests := []Conflict{
			{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(10)}}},
		}
		result := ResolveConflicts(bracket, requests, players)
		_, count := checkConflict(bracket, requests, result)
		if count != 0 {
			t.Errorf("run %d: request unfulfilled", i)
		}
	}
}

// TestResolveHighPriorityRequestBeatsLowPriorityConflict verifies that a
// priority-3 request is honored even when it creates a priority-1 conflict
// as a side effect. The TO's explicit "force this match" should outweigh a
// soft "prefer to avoid this other match".
func TestResolveHighPriorityRequestBeatsLowPriorityConflict(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// Request P1 vs P6 (priority -3, must happen).
	// Conflict P2 vs P7 (priority 1, prefer to avoid) — which they currently do.
	cons := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
		{Priority: 1, Players: []Player{{ID: int64(2)}, {ID: int64(7)}}},
	}

	result := ResolveConflicts(bracket, cons, players)

	// The high-priority request must be fulfilled.
	requestOnly := []Conflict{cons[0]}
	_, count := checkConflict(bracket, requestOnly, result)
	if count != 0 {
		t.Error("high-priority set request should be fulfilled even at the cost of a low-priority conflict")
	}
}

// TestResolveMultipleRequestsAtScale combines multiple requests at 32 players
// to verify the algorithm handles load.
func TestResolveMultipleRequestsAtScale(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 32-player multi-request stress test in short mode")
	}
	resolvesRequest(t, 32, []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(15)}}},
		{Priority: -3, Players: []Player{{ID: int64(3)}, {ID: int64(20)}}},
		{Priority: -3, Players: []Player{{ID: int64(5)}, {ID: int64(25)}}},
	})
}

// TestResolveImpossibleRequestDoesNotCrash verifies that if a request cannot
// be fulfilled (e.g. the same player paired with themselves — the checker
// ignores it), ResolveConflicts still terminates without panicking.
func TestResolveImpossibleRequestDoesNotCrash(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// A degenerate request: same player twice. check() requires TWO distinct
	// players to trigger, so this is always unfulfilled.
	cons := []Conflict{{
		Priority: -3,
		Players:  []Player{{ID: int64(1)}, {ID: int64(1)}},
	}}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ResolveConflicts panicked on impossible request: %v", r)
		}
	}()

	_ = ResolveConflicts(bracket, cons, players)
}
