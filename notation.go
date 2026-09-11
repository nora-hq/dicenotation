// Package dicenotation parses and rolls dice expressions like "3d6+2" or
// "4d6dl1", the shorthand used in tabletop games to describe a set of dice
// and how to combine the results into a single number.
package dicenotation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// pattern matches: optional count, "d", sides, optional keep/drop
// modifier, optional +/- modifier. Examples: "d20", "3d6", "4d6dl1",
// "2d20kh1", "1d8+3".
var pattern = regexp.MustCompile(`^(\d*)d(\d+)(kh|kl|dh|dl)?(\d*)([+-]\d+)?$`)

// Notation is a parsed dice expression, e.g. "4d6dl1" (roll four six-sided
// dice, drop the lowest one).
type Notation struct {
	Count       int // number of dice rolled
	Sides       int // faces per die
	DropLowest  int // how many of the lowest rolls to discard
	DropHighest int // how many of the highest rolls to discard
	Modifier    int // flat amount added to the kept total
}

// Parse reads a dice expression and returns the Notation it describes.
// A bare count defaults to 1 ("d20" is the same as "1d20").
func Parse(expr string) (Notation, error) {
	s := strings.ToLower(strings.TrimSpace(expr))
	m := pattern.FindStringSubmatch(s)
	if m == nil {
		return Notation{}, fmt.Errorf("dicenotation: invalid expression %q", expr)
	}

	n := Notation{Count: 1, Sides: 1}

	if m[1] != "" {
		n.Count, _ = strconv.Atoi(m[1])
	}
	n.Sides, _ = strconv.Atoi(m[2])

	if n.Count < 1 {
		return Notation{}, fmt.Errorf("dicenotation: %q needs at least one die", expr)
	}
	if n.Sides < 1 {
		return Notation{}, fmt.Errorf("dicenotation: %q needs at least one side", expr)
	}

	if kind := m[3]; kind != "" {
		amount := 1
		if m[4] != "" {
			amount, _ = strconv.Atoi(m[4])
		}
		switch kind {
		case "dl":
			n.DropLowest = amount
		case "dh":
			n.DropHighest = amount
		case "kh":
			// keep the N highest means dropping the rest as lowest
			n.DropLowest = n.Count - amount
		case "kl":
			// keep the N lowest means dropping the rest as highest
			n.DropHighest = n.Count - amount
		}
	}

	if n.DropLowest+n.DropHighest >= n.Count {
		return Notation{}, fmt.Errorf("dicenotation: %q drops all the dice, nothing left to keep", expr)
	}
	if n.DropLowest < 0 || n.DropHighest < 0 {
		return Notation{}, fmt.Errorf("dicenotation: %q keeps more dice than it rolls", expr)
	}

	if m[5] != "" {
		n.Modifier, _ = strconv.Atoi(m[5])
	}

	return n, nil
}

// String renders the Notation back into dice notation.
func (n Notation) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%dd%d", n.Count, n.Sides)
	if n.DropLowest > 0 {
		fmt.Fprintf(&b, "dl%d", n.DropLowest)
	}
	if n.DropHighest > 0 {
		fmt.Fprintf(&b, "dh%d", n.DropHighest)
	}
	if n.Modifier != 0 {
		fmt.Fprintf(&b, "%+d", n.Modifier)
	}
	return b.String()
}
