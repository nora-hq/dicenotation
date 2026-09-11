# dicenotation

A small Go library for parsing and rolling dice notation, the shorthand used
in tabletop games to describe how many dice to roll and how to turn the
results into a number: `3d6`, `1d20+5`, `4d6dl1`.

There's no CLI here, just a package. If you're building a character sheet
tool, a game bot, or a random-encounter generator, you probably want to
parse a string a user typed and turn it into a number (or the individual
rolls, if you want to show your work). That's what this does.

## Install

```
go get github.com/nora-hq/dicenotation
```

## Usage

Roll an expression directly:

```go
res, err := dicenotation.Roll("4d6dl1")
if err != nil {
    log.Fatal(err)
}
fmt.Println(res) // [5 2 6 1] (dropped 1) = 13
fmt.Println(res.Total)
```

Parse once, roll many times (useful if you're rolling the same expression
repeatedly, or want to control the random source):

```go
n, err := dicenotation.Parse("2d20kh1") // roll two d20s, keep the higher
if err != nil {
    log.Fatal(err)
}

r := rand.New(rand.NewSource(42))
for i := 0; i < 5; i++ {
    fmt.Println(n.Roll(r))
}
```

## Notation

| Syntax     | Meaning                                    |
|------------|---------------------------------------------|
| `d20`      | roll one twenty-sided die                   |
| `3d6`      | roll three six-sided dice, sum them         |
| `1d8+3`    | roll one d8, add 3 to the result             |
| `4d6dl1`   | roll four d6, drop the lowest one           |
| `4d6dh1`   | roll four d6, drop the highest one          |
| `2d20kh1`  | roll two d20, keep the highest one          |
| `2d20kl1`  | roll two d20, keep the lowest one           |

A bare count defaults to 1, so `d20` and `1d20` are the same expression.

## What this doesn't do (yet)

No exploding dice, no rerolls, no compound expressions like `2d6+1d4`. See
the issue tracker for what's planned.

## License

MIT, see [LICENSE](LICENSE).
