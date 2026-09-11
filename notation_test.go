package dicenotation

import (
	"math/rand"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		expr string
		want Notation
	}{
		{"d20", Notation{Count: 1, Sides: 20}},
		{"3d6", Notation{Count: 3, Sides: 6}},
		{"1d8+3", Notation{Count: 1, Sides: 8, Modifier: 3}},
		{"4d6dl1", Notation{Count: 4, Sides: 6, DropLowest: 1}},
		{"2d20kh1", Notation{Count: 2, Sides: 20, DropLowest: 1}},
		{"2d20kl1", Notation{Count: 2, Sides: 20, DropHighest: 1}},
	}
	for _, c := range cases {
		got, err := Parse(c.expr)
		if err != nil {
			t.Fatalf("Parse(%q) returned error: %v", c.expr, err)
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.expr, got, c.want)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	for _, expr := range []string{"", "d0", "0d6", "4d6dl4", "abc"} {
		if _, err := Parse(expr); err == nil {
			t.Errorf("Parse(%q) expected an error, got nil", expr)
		}
	}
}

func TestRollBounds(t *testing.T) {
	n, err := Parse("4d6dl1")
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		res := n.Roll(r)
		if len(res.Rolls) != 4 {
			t.Fatalf("got %d rolls, want 4", len(res.Rolls))
		}
		if len(res.Dropped) != 1 {
			t.Fatalf("got %d dropped, want 1", len(res.Dropped))
		}
		if res.Total < 3 || res.Total > 18 {
			t.Fatalf("total %d out of range [3, 18]", res.Total)
		}
	}
}
