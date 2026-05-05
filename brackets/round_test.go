package brackets

import "testing"

func TestParseRound(t *testing.T) {
	tests := []struct {
		input   string
		want    Round
		wantErr bool
	}{
		{"WR0", Round{Number: 0}, false},
		{"WR2", Round{Number: 2}, false},
		{"LR1", Round{Losers: true, Number: 1}, false},
		{"lr3", Round{Losers: true, Number: 3}, false},
		{"wr1", Round{Number: 1}, false},
		{"", Round{}, true},
		{"R1", Round{}, true},
		{"WR", Round{}, true},
		{"WRabc", Round{}, true},
	}

	for _, tt := range tests {
		got, err := ParseRound(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseRound(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseRound(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestRoundString(t *testing.T) {
	if s := (Round{Number: 2}).String(); s != "WR2" {
		t.Errorf("Expected WR2, got %s", s)
	}
	if s := (Round{Losers: true, Number: 1}).String(); s != "LR1" {
		t.Errorf("Expected LR1, got %s", s)
	}
}

func TestRoundFromStartGG(t *testing.T) {
	r := RoundFromStartGG(3)
	if r.Losers || r.Number != 3 {
		t.Errorf("Expected WR3, got %v", r)
	}

	r = RoundFromStartGG(-2)
	if !r.Losers || r.Number != 2 {
		t.Errorf("Expected LR2, got %v", r)
	}
}
