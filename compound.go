package dicenotation

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// signedToken splits a compound expression into its signed pieces:
// "2d6+1d4-3" becomes the tokens "+2d6", "+1d4", "-3".
var signedToken = regexp.MustCompile(`[+-][^+-]*`)

// Term is one signed dice notation within a compound Expression.
type Term struct {
	Notation Notation
	Negative bool
}

// Expression is a sum of dice terms and a flat modifier, e.g. "2d6+1d4-3"
// (roll two six-sided dice, add one four-sided die, then subtract three).
// A plain Notation like "3d6+2" is a special case of an Expression with a
// single term.
type Expression struct {
	Terms    []Term
	Modifier int
}

// ParseExpression reads a compound dice expression and returns the
// Expression it describes. Each term is parsed the same way as Parse, so
// per-term keep/drop modifiers like "4d6dl1+2d8" are still allowed; only
// flat modifiers (bare numbers) are pulled out and summed separately.
func ParseExpression(expr string) (Expression, error) {
	s := strings.ToLower(strings.TrimSpace(expr))
	if s == "" {
		return Expression{}, fmt.Errorf("dicenotation: invalid expression %q", expr)
	}
	if s[0] != '+' && s[0] != '-' {
		s = "+" + s
	}

	tokens := signedToken.FindAllString(s, -1)
	if strings.Join(tokens, "") != s {
		return Expression{}, fmt.Errorf("dicenotation: invalid expression %q", expr)
	}

	var e Expression
	for _, tok := range tokens {
		negative := tok[0] == '-'
		body := tok[1:]
		if body == "" {
			return Expression{}, fmt.Errorf("dicenotation: invalid expression %q", expr)
		}

		if n, err := Parse(body); err == nil {
			e.Terms = append(e.Terms, Term{Notation: n, Negative: negative})
			continue
		}

		amount, err := strconv.Atoi(body)
		if err != nil {
			return Expression{}, fmt.Errorf("dicenotation: invalid expression %q", expr)
		}
		if negative {
			amount = -amount
		}
		e.Modifier += amount
	}

	if len(e.Terms) == 0 {
		return Expression{}, fmt.Errorf("dicenotation: %q has no dice", expr)
	}

	return e, nil
}

// String renders the Expression back into dice notation.
func (e Expression) String() string {
	var b strings.Builder
	for i, t := range e.Terms {
		if t.Negative {
			b.WriteString("-")
		} else if i > 0 {
			b.WriteString("+")
		}
		b.WriteString(t.Notation.String())
	}
	if e.Modifier != 0 {
		fmt.Fprintf(&b, "%+d", e.Modifier)
	}
	return b.String()
}

// CompoundResult is the outcome of rolling an Expression.
type CompoundResult struct {
	Expression Expression
	Results    []Result // one per term, in the order the terms appear
	Total      int
}

// String renders each term's rolls and the grand total, e.g.
// "[3 5] + [2] = 10".
func (r CompoundResult) String() string {
	parts := make([]string, len(r.Results))
	for i, res := range r.Results {
		s := rollsString(res.Rolls, res.Dropped)
		switch {
		case r.Expression.Terms[i].Negative:
			s = "- " + s
		case i > 0:
			s = "+ " + s
		}
		parts[i] = s
	}
	out := strings.Join(parts, " ")
	if r.Expression.Modifier != 0 {
		out += fmt.Sprintf(" %+d", r.Expression.Modifier)
	}
	return fmt.Sprintf("%s = %d", out, r.Total)
}

// Roll rolls every term in the expression using r as the source of
// randomness, subtracting the total of any term marked Negative, and sums
// the results with the flat modifier. Passing nil uses the package's
// default random source.
func (e Expression) Roll(r *rand.Rand) CompoundResult {
	total := e.Modifier
	results := make([]Result, len(e.Terms))
	for i, t := range e.Terms {
		res := t.Notation.Roll(r)
		results[i] = res
		if t.Negative {
			total -= res.Total
		} else {
			total += res.Total
		}
	}
	return CompoundResult{Expression: e, Results: results, Total: total}
}

// RollExpression parses expr and rolls it in one step, using the package's
// default random source.
func RollExpression(expr string) (CompoundResult, error) {
	e, err := ParseExpression(expr)
	if err != nil {
		return CompoundResult{}, err
	}
	return e.Roll(nil), nil
}
