package dicenotation

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// Result is the outcome of rolling a Notation.
type Result struct {
	Notation Notation
	Rolls    []int // every die, in the order it was rolled
	Dropped  []int // values discarded by a keep/drop modifier
	Total    int   // sum of the kept dice plus the modifier
}

// String renders the individual rolls and the total, e.g. "[5 2 6] -3 = 10".
func (r Result) String() string {
	s := rollsString(r.Rolls, r.Dropped)
	if r.Notation.Modifier != 0 {
		s += fmt.Sprintf(" %+d", r.Notation.Modifier)
	}
	return fmt.Sprintf("%s = %d", s, r.Total)
}

// rollsString renders the dice values and any dropped dice, e.g.
// "[5 2 6 1] (dropped 1)". It's shared by Result and CompoundResult so a
// compound expression's per-term output matches a plain Result's.
func rollsString(rolls, dropped []int) string {
	parts := make([]string, len(rolls))
	for i, v := range rolls {
		parts[i] = fmt.Sprintf("%d", v)
	}
	s := fmt.Sprintf("[%s]", strings.Join(parts, " "))
	if len(dropped) > 0 {
		d := make([]string, len(dropped))
		for i, v := range dropped {
			d[i] = fmt.Sprintf("%d", v)
		}
		s += fmt.Sprintf(" (dropped %s)", strings.Join(d, " "))
	}
	return s
}

// Roll rolls the dice described by n using r as the source of randomness.
// Passing nil uses the package's default source, which is seeded from the
// runtime and safe for concurrent use.
func (n Notation) Roll(r *rand.Rand) Result {
	intn := rand.Intn
	if r != nil {
		intn = r.Intn
	}

	rolls := make([]int, n.Count)
	for i := range rolls {
		rolls[i] = intn(n.Sides) + 1
	}

	kept := append([]int(nil), rolls...)
	sort.Ints(kept)

	var dropped []int
	if n.DropLowest > 0 {
		dropped = append(dropped, kept[:n.DropLowest]...)
		kept = kept[n.DropLowest:]
	}
	if n.DropHighest > 0 {
		end := len(kept) - n.DropHighest
		dropped = append(dropped, kept[end:]...)
		kept = kept[:end]
	}

	total := n.Modifier
	for _, v := range kept {
		total += v
	}

	return Result{Notation: n, Rolls: rolls, Dropped: dropped, Total: total}
}

// Roll parses expr and rolls it in one step, using the package's default
// random source. It is a convenience for callers who don't need to reuse
// a parsed Notation or supply their own *rand.Rand.
func Roll(expr string) (Result, error) {
	n, err := Parse(expr)
	if err != nil {
		return Result{}, err
	}
	return n.Roll(nil), nil
}
