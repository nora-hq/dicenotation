package dicenotation

import (
	"math/rand"
	"testing"
)

func TestParseExpression(t *testing.T) {
	cases := []struct {
		expr string
		want Expression
	}{
		{"3d6", Expression{Terms: []Term{{Notation: Notation{Count: 3, Sides: 6}}}}},
		{"1d8+3", Expression{Terms: []Term{{Notation: Notation{Count: 1, Sides: 8}}}, Modifier: 3}},
		{
			"2d6+1d4",
			Expression{Terms: []Term{
				{Notation: Notation{Count: 2, Sides: 6}},
				{Notation: Notation{Count: 1, Sides: 4}},
			}},
		},
		{
			"4d6dl1+2d8-2",
			Expression{Terms: []Term{
				{Notation: Notation{Count: 4, Sides: 6, DropLowest: 1}},
				{Notation: Notation{Count: 2, Sides: 8}},
			}, Modifier: -2},
		},
		{
			"-2d6+1d4",
			Expression{Terms: []Term{
				{Notation: Notation{Count: 2, Sides: 6}, Negative: true},
				{Notation: Notation{Count: 1, Sides: 4}},
			}},
		},
	}
	for _, c := range cases {
		got, err := ParseExpression(c.expr)
		if err != nil {
			t.Fatalf("ParseExpression(%q) returned error: %v", c.expr, err)
		}
		if len(got.Terms) != len(c.want.Terms) || got.Modifier != c.want.Modifier {
			t.Fatalf("ParseExpression(%q) = %+v, want %+v", c.expr, got, c.want)
		}
		for i, term := range got.Terms {
			if term != c.want.Terms[i] {
				t.Errorf("ParseExpression(%q) term %d = %+v, want %+v", c.expr, i, term, c.want.Terms[i])
			}
		}
	}
}

func TestParseExpressionInvalid(t *testing.T) {
	for _, expr := range []string{"", "+", "2d6+", "2d6++1d4", "5", "d6+", "2d6+abc"} {
		if _, err := ParseExpression(expr); err == nil {
			t.Errorf("ParseExpression(%q) expected an error, got nil", expr)
		}
	}
}

func TestExpressionRollBounds(t *testing.T) {
	e, err := ParseExpression("2d6+1d4-1")
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		res := e.Roll(r)
		if len(res.Results) != 2 {
			t.Fatalf("got %d term results, want 2", len(res.Results))
		}
		if res.Total < 2 || res.Total > 15 {
			t.Fatalf("total %d out of range [2, 15]", res.Total)
		}
	}
}

func TestExpressionRollNegativeTerm(t *testing.T) {
	e, err := ParseExpression("2d6-1d4")
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		res := e.Roll(r)
		if res.Total < 2-4 || res.Total > 12-1 {
			t.Fatalf("total %d out of range [-2, 11]", res.Total)
		}
	}
}
