package conflicts

import (
	"testing"

	"github.com/AaronLieb/octagon/brackets"
)

func TestConflictCheck(t *testing.T) {
	conflict := Conflict{
		Priority: 1,
		Reason:   "test conflict",
		Players: []Player{
			{Name: "Player1", ID: int64(123)},
			{Name: "Player2", ID: int64(456)},
		},
	}

	// Test matching conflict
	if !conflict.check(int64(123), int64(456)) {
		t.Error("Expected conflict to match players 123 and 456")
	}

	// Test reversed order still matches
	if !conflict.check(int64(456), int64(123)) {
		t.Error("Expected conflict to match players in reversed order")
	}

	// Test non-matching conflict
	if conflict.check(int64(123), int64(789)) {
		t.Error("Expected no conflict between players 123 and 789")
	}

	// Test same player
	if conflict.check(int64(123), int64(123)) {
		t.Error("Expected no conflict for same player")
	}
}

func TestCheckConflict(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(123), Name: "Player1"},
		{ID: int64(456), Name: "Player2"},
		{ID: int64(789), Name: "Player3"},
	}

	conflicts := []Conflict{
		{
			Priority: 1,
			Players: []Player{
				{ID: int64(123)}, {ID: int64(456)},
			},
		},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2}, // Players 123 vs 456 - should conflict
			{Player1: 2, Player2: 3}, // Players 456 vs 789 - no conflict
		},
	}

	score, count := checkConflict(bracket, conflicts, players)

	if count != 1 {
		t.Errorf("Expected 1 conflict, got %d", count)
	}

	if score != 3.0 { // 2 + priority(1)
		t.Errorf("Expected score 3.0, got %.1f", score)
	}
}

func TestCheckConflictNoMatches(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
		{ID: int64(3), Name: "P3"},
	}

	conflicts := []Conflict{
		{Priority: 1, Players: []Player{{ID: int64(1)}, {ID: int64(3)}}},
	}

	// Players 1 and 3 never play each other
	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
			{Player1: 2, Player2: 3},
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	if score != 0.0 {
		t.Errorf("Expected score 0, got %.1f", score)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestCheckConflictMultipleOnSameSet(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
	}

	conflicts := []Conflict{
		{Priority: 1, Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
		{Priority: 2, Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	if count != 2 {
		t.Errorf("Expected 2 conflicts, got %d", count)
	}
	// (2+1) + (2+2) = 7
	if score != 7.0 {
		t.Errorf("Expected score 7.0, got %.1f", score)
	}
}

func TestListUnresolvedConflicts(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
		{ID: int64(3), Name: "P3"},
	}

	conflicts := []Conflict{
		{Priority: 1, Reason: "unresolved", Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
		{Priority: 1, Reason: "resolved", Players: []Player{{ID: int64(1)}, {ID: int64(3)}}},
	}

	// Only P1 vs P2 play each other
	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
		},
	}

	unresolved := listUnresolvedConflicts(bracket, conflicts, players)
	if len(unresolved) != 1 {
		t.Fatalf("Expected 1 unresolved conflict, got %d", len(unresolved))
	}
	if unresolved[0].Reason != "unresolved" {
		t.Errorf("Expected unresolved conflict, got %s", unresolved[0].Reason)
	}
}

func TestIsRequest(t *testing.T) {
	pos := Conflict{Priority: 1}
	neg := Conflict{Priority: -1}
	zero := Conflict{Priority: 0}

	if pos.isRequest() {
		t.Error("Positive priority should not be a request")
	}
	if !neg.isRequest() {
		t.Error("Negative priority should be a request")
	}
	if zero.isRequest() {
		t.Error("Zero priority should not be a request")
	}
}

func TestSetRequestPenaltyWhenNotFound(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
		{ID: int64(3), Name: "P3"},
	}

	// Request that P1 and P3 play each other
	conflicts := []Conflict{
		{Priority: -2, Players: []Player{{ID: int64(1)}, {ID: int64(3)}}},
	}

	// P1 vs P3 never happens
	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
			{Player1: 2, Player2: 3},
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	if count != 1 {
		t.Errorf("Expected 1 unfulfilled request, got %d", count)
	}
	// 2 + abs(-2) = 4
	if score != 4.0 {
		t.Errorf("Expected score 4.0, got %.1f", score)
	}
}

func TestSetRequestNoPenaltyWhenFulfilled(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
	}

	// Request that P1 and P2 play each other
	conflicts := []Conflict{
		{Priority: -1, Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
	}

	// P1 vs P2 does happen
	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	if count != 0 {
		t.Errorf("Expected 0 unfulfilled requests, got %d", count)
	}
	if score != 0.0 {
		t.Errorf("Expected score 0, got %.1f", score)
	}
}

func TestMixedConflictsAndRequests(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
		{ID: int64(3), Name: "P3"},
	}

	conflicts := []Conflict{
		// Conflict: P1 vs P2 should NOT play (but they do)
		{Priority: 1, Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
		// Request: P1 vs P3 SHOULD play (but they don't)
		{Priority: -2, Players: []Player{{ID: int64(1)}, {ID: int64(3)}}},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2}, // P1 vs P2 happens
			{Player1: 2, Player2: 3}, // P2 vs P3 happens
		},
	}

	score, count := checkConflict(bracket, conflicts, players)
	if count != 2 {
		t.Errorf("Expected 2 issues (1 conflict + 1 unfulfilled request), got %d", count)
	}
	// Conflict: 2+1=3, Request unfulfilled: 2+abs(-2)=4, total=7
	if score != 7.0 {
		t.Errorf("Expected score 7.0, got %.1f", score)
	}
}

func TestListUnresolvedWithRequests(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1), Name: "P1"},
		{ID: int64(2), Name: "P2"},
		{ID: int64(3), Name: "P3"},
	}

	conflicts := []Conflict{
		// Conflict: P1 vs P2 - they DO play, so this is unresolved
		{Priority: 1, Reason: "conflict-unresolved", Players: []Player{{ID: int64(1)}, {ID: int64(2)}}},
		// Request: P1 vs P3 - they DON'T play, so this is unresolved
		{Priority: -1, Reason: "request-unfulfilled", Players: []Player{{ID: int64(1)}, {ID: int64(3)}}},
		// Request: P2 vs P3 - they DO play, so this is fulfilled (resolved)
		{Priority: -1, Reason: "request-fulfilled", Players: []Player{{ID: int64(2)}, {ID: int64(3)}}},
	}

	bracket := &brackets.Bracket{
		Sets: []*brackets.Set{
			{Player1: 1, Player2: 2},
			{Player1: 2, Player2: 3},
		},
	}

	unresolved := listUnresolvedConflicts(bracket, conflicts, players)
	if len(unresolved) != 2 {
		t.Fatalf("Expected 2 unresolved, got %d", len(unresolved))
	}

	reasons := map[string]bool{}
	for _, u := range unresolved {
		reasons[u.Reason] = true
	}
	if !reasons["conflict-unresolved"] {
		t.Error("Expected conflict-unresolved to be listed")
	}
	if !reasons["request-unfulfilled"] {
		t.Error("Expected request-unfulfilled to be listed")
	}
	if reasons["request-fulfilled"] {
		t.Error("Fulfilled request should not be listed as unresolved")
	}
}

func TestResolveConflictsWithSetRequest(t *testing.T) {
	// In an 8-player bracket, seed 1 plays seeds 8, 4, 2
	// P1 (seed 1) and P6 (seed 6) never meet naturally
	// The optimizer should rearrange seeds so they meet
	players := []brackets.Player{
		{ID: int64(1), Name: "P1", Rating: 8},
		{ID: int64(2), Name: "P2", Rating: 7},
		{ID: int64(3), Name: "P3", Rating: 6},
		{ID: int64(4), Name: "P4", Rating: 5},
		{ID: int64(5), Name: "P5", Rating: 4},
		{ID: int64(6), Name: "P6", Rating: 3},
		{ID: int64(7), Name: "P7", Rating: 2},
		{ID: int64(8), Name: "P8", Rating: 1},
	}

	bracket := brackets.CreateBracket(8)

	// Verify P1 vs P6 don't meet in default seeding
	initialRequests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(6)}}},
	}
	_, initialCount := checkConflict(bracket, initialRequests, players)
	if initialCount != 0 {
		t.Log("Confirmed: P1 and P6 don't meet in default seeding")
	} else {
		t.Fatal("P1 and P6 should not meet in default 8-player bracket seeding")
	}

	result := ResolveConflicts(bracket, initialRequests, players)

	// Verify the request is fulfilled after resolution
	_, count := checkConflict(bracket, initialRequests, result)
	if count != 0 {
		t.Error("Expected set request to be fulfilled after resolution")
	}
}

func wr(n int) *brackets.Round { return &brackets.Round{Number: n} }
func lr(n int) *brackets.Round { return &brackets.Round{Losers: true, Number: n} }

func TestConflictWithRoundOnlyPenalizesMatchingRound(t *testing.T) {
	// 8-player bracket: WR1 has seeds 1v8, 4v5, 2v7, 3v6
	players := []brackets.Player{
		{ID: int64(1)}, {ID: int64(2)}, {ID: int64(3)}, {ID: int64(4)},
		{ID: int64(5)}, {ID: int64(6)}, {ID: int64(7)}, {ID: int64(8)},
	}
	bracket := brackets.CreateBracket(8)

	// Conflict on WR1: P1 vs P8 play in WR1, should penalize
	cons := []Conflict{
		{Priority: 1, Round: wr(1), Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
	}
	score, count := checkConflict(bracket, cons, players)
	if count != 1 || score != 3.0 {
		t.Errorf("Expected penalty for WR1 conflict, got score=%.1f count=%d", score, count)
	}

	// Same conflict but on WR2: P1 vs P8 don't play in WR2, no penalty
	cons[0].Round = wr(2)
	score, count = checkConflict(bracket, cons, players)
	if count != 0 || score != 0.0 {
		t.Errorf("Expected no penalty for wrong round, got score=%.1f count=%d", score, count)
	}
}

func TestSetRequestWithRoundOnlyChecksMatchingRound(t *testing.T) {
	// 8-player bracket: WR1 has 1v8, 4v5, 2v7, 3v6
	// P1 and P4 meet in WR3 (winners semis) but NOT WR1
	players := []brackets.Player{
		{ID: int64(1)}, {ID: int64(2)}, {ID: int64(3)}, {ID: int64(4)},
		{ID: int64(5)}, {ID: int64(6)}, {ID: int64(7)}, {ID: int64(8)},
	}
	bracket := brackets.CreateBracket(8)

	// Request P1 vs P4 in WR1: they don't play in WR1, should penalize
	cons := []Conflict{
		{Priority: -2, Round: wr(1), Players: []Player{{ID: int64(1)}, {ID: int64(4)}}},
	}
	score, count := checkConflict(bracket, cons, players)
	if count != 1 {
		t.Errorf("Expected unfulfilled request for WR1, got count=%d", count)
	}
	if score != 4.0 { // 2 + abs(-2)
		t.Errorf("Expected score 4.0, got %.1f", score)
	}

	// Request P1 vs P8 in WR1: they DO play in WR1, no penalty
	cons = []Conflict{
		{Priority: -2, Round: wr(1), Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
	}
	score, count = checkConflict(bracket, cons, players)
	if count != 0 || score != 0.0 {
		t.Errorf("Expected fulfilled request, got score=%.1f count=%d", score, count)
	}
}

func TestConflictWithLosersRound(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1)}, {ID: int64(2)}, {ID: int64(3)}, {ID: int64(4)},
		{ID: int64(5)}, {ID: int64(6)}, {ID: int64(7)}, {ID: int64(8)},
	}
	bracket := brackets.CreateBracket(8)

	// LR1 conflict should not match WR1 sets
	cons := []Conflict{
		{Priority: 1, Round: lr(1), Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
	}
	score, count := checkConflict(bracket, cons, players)
	// P1 vs P8 is in WR1, not LR1, so no penalty
	if count != 0 || score != 0.0 {
		t.Errorf("Expected no penalty for LR1 when match is in WR1, got score=%.1f count=%d", score, count)
	}
}

func TestListUnresolvedWithRound(t *testing.T) {
	players := []brackets.Player{
		{ID: int64(1)}, {ID: int64(2)}, {ID: int64(3)}, {ID: int64(4)},
		{ID: int64(5)}, {ID: int64(6)}, {ID: int64(7)}, {ID: int64(8)},
	}
	bracket := brackets.CreateBracket(8)

	cons := []Conflict{
		// Conflict: P1 vs P8 in WR1 — they play, so unresolved
		{Priority: 1, Round: wr(1), Reason: "wr1-conflict", Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
		// Conflict: P1 vs P8 in WR2 — they don't play in WR2, so resolved
		{Priority: 1, Round: wr(2), Reason: "wr2-conflict", Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
		// Request: P1 vs P4 in WR1 — they don't play in WR1, so unresolved
		{Priority: -1, Round: wr(1), Reason: "wr1-request", Players: []Player{{ID: int64(1)}, {ID: int64(4)}}},
	}

	unresolved := listUnresolvedConflicts(bracket, cons, players)
	if len(unresolved) != 2 {
		t.Fatalf("Expected 2 unresolved, got %d", len(unresolved))
	}

	reasons := map[string]bool{}
	for _, u := range unresolved {
		reasons[u.Reason] = true
	}
	if !reasons["wr1-conflict"] {
		t.Error("Expected wr1-conflict to be unresolved")
	}
	if !reasons["wr1-request"] {
		t.Error("Expected wr1-request to be unresolved")
	}
	if reasons["wr2-conflict"] {
		t.Error("wr2-conflict should be resolved (wrong round)")
	}
}
