package conflicts

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AaronLieb/octagon/brackets"
	"github.com/charmbracelet/log"
)

// captureLog redirects the default logger output to a buffer for the
// duration of fn, returning whatever was written.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()

	buf := &bytes.Buffer{}
	orig := log.Default()
	defer log.SetDefault(orig)

	l := log.New(buf)
	l.SetLevel(log.InfoLevel)
	log.SetDefault(l)

	fn()
	return buf.String()
}

func TestPrintSetRequests_EmptyIsQuiet(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	out := captureLog(t, func() {
		printSetRequests(bracket, nil, players)
	})

	if out != "" {
		t.Errorf("expected no output for empty request list, got:\n%s", out)
	}
}

func TestPrintSetRequests_FiltersNonRequests(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// A non-request conflict should be filtered out defensively.
	cons := []Conflict{{Priority: 3, Players: []Player{{ID: int64(1)}, {ID: int64(5)}}}}

	out := captureLog(t, func() {
		printSetRequests(bracket, cons, players)
	})

	if strings.Contains(out, "set request") {
		t.Errorf("non-request conflicts should not produce set-request output, got:\n%s", out)
	}
}

func TestPrintSetRequests_ListsEachRequest(t *testing.T) {
	players := makePlayers(16)
	names := []string{
		"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
		"india", "juliet", "kilo", "lima", "mike", "november", "oscar", "papa",
	}
	for i := range players {
		players[i].Name = names[i]
	}
	bracket := brackets.CreateBracket(16)

	requests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(10)}}},
		{Priority: -1, Round: wr(1), Players: []Player{{ID: int64(3)}, {ID: int64(14)}}},
	}

	out := captureLog(t, func() {
		printSetRequests(bracket, requests, players)
	})

	if !strings.Contains(out, "2 unfulfilled set request(s)") {
		t.Errorf("missing header, output:\n%s", out)
	}

	for _, want := range []string{"alpha", "juliet", "charlie", "november"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected player %q in output:\n%s", want, out)
		}
	}

	if !strings.Contains(out, "r3") || !strings.Contains(out, "r1") {
		t.Errorf("expected priority labels r3 and r1 in output:\n%s", out)
	}

	if !strings.Contains(out, "any round") {
		t.Errorf("expected 'any round' for request without round:\n%s", out)
	}
	if !strings.Contains(out, "WR1") {
		t.Errorf("expected 'WR1' for round-specific request:\n%s", out)
	}
}

func TestPrintSetRequests_ShowsCurrentRound(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// In an 8-player bracket, seed 1 and seed 8 meet in WR1. A request for
	// that pair at "any round" would normally be fulfilled, but this
	// function is called only on unresolved requests — we verify the
	// "current round" column still reports where the pair would meet.
	requests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},
	}

	out := captureLog(t, func() {
		printSetRequests(bracket, requests, players)
	})

	if !strings.Contains(out, "WR1") {
		t.Errorf("expected current round 'WR1' in output:\n%s", out)
	}
}

func TestPrintSetRequests_ShowsDashWhenPairDoesNotMeet(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// A request for a non-existent player id: the pair doesn't meet
	// anywhere in the bracket.
	requests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}, {ID: int64(999)}}},
	}

	out := captureLog(t, func() {
		printSetRequests(bracket, requests, players)
	})

	if !strings.Contains(out, "—") {
		t.Errorf("expected em-dash for no current round, output:\n%s", out)
	}
}

func TestPrintSetRequests_SkipsMalformedRequests(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// One malformed (single player), one valid.
	requests := []Conflict{
		{Priority: -3, Players: []Player{{ID: int64(1)}}},
		{Priority: -3, Players: []Player{{ID: int64(2)}, {ID: int64(7)}}},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("printSetRequests panicked: %v", r)
		}
	}()

	out := captureLog(t, func() {
		printSetRequests(bracket, requests, players)
	})

	// Header must reflect the printable count (1), not the input count (2).
	if !strings.Contains(out, "1 unfulfilled set request(s)") {
		t.Errorf("expected header to count only printable requests, output:\n%s", out)
	}
}

func TestPrintUnresolved_NoDoubleCounting(t *testing.T) {
	players := makePlayers(8)
	bracket := brackets.CreateBracket(8)

	// Mix: one unresolvable regular conflict and one unfulfilled set request.
	cons := []Conflict{
		{Priority: 3, Players: []Player{{ID: int64(1)}, {ID: int64(8)}}},                   // P1 vs P8 meet in WR1, so this regular conflict is unresolved
		{Priority: -3, Round: wr(1), Players: []Player{{ID: int64(2)}, {ID: int64(3)}}},    // P2 and P3 are in different WR1 sets, so this request is unresolved
	}

	out := captureLog(t, func() {
		printUnresolved(bracket, cons, players)
	})

	// The set request's priority label (r3) appears once.
	if count := strings.Count(out, "r3"); count != 1 {
		t.Errorf("set request should appear exactly once, got %d occurrences:\n%s", count, out)
	}
	// The regular conflict's priority label (p3) appears once.
	if count := strings.Count(out, "p3"); count != 1 {
		t.Errorf("regular conflict should appear exactly once, got %d occurrences:\n%s", count, out)
	}
}
