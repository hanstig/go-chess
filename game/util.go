package game

import (
	"bytes"
	"fmt"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

func isStraightMove(m Move) bool {
	dx, dy := m.deltas()
	return (dx == 0) != (dy == 0)
}

func isDiagonalMove(m Move) bool {
	dx, dy := m.deltas()
	return abs(dx) == abs(dy) && dx != 0
}

func isPathClear(g *Game, m Move) bool {
	dx, dy := m.deltas()
	dx, dy = sign(dx), sign(dy)

	x, y := m.From.X+dx, m.From.Y+dy

	for x != m.To.X || y != m.To.Y {
		if g.Board[y][x] != nil {
			return false
		}

		x += dx
		y += dy
	}

	return true
}

//
// Debug functions
//

func PrintMoves(g *Game, mvs []Move) {
	printchar := func(s string, clr string, b *bytes.Buffer) {
		clear := "\033[0m"
		beg, end := clear, clear
		switch clr {
		case "red":
			beg = "\033[41m"
		case "green":
			beg = "\033[42m"
		}
		fmt.Fprintf(b, "%s%s%s", beg, s, end)
	}

	var buf bytes.Buffer

	for _, m := range mvs {
		buf.Reset()
		clr := ""

		for y := range 8 {
			for x := range 8 {
				s := Square{X: x, Y: y}
				if s == m.To {
					clr = "red"
				} else if s == m.From {
					clr = "green"
				} else {
					clr = ""
				}

				p := g.Board[y][x]

				if p == nil {
					printchar(".", clr, &buf)
				} else {
					c := string(p.Type)
					if p.Color == Black {
						c = strings.ToLower(c)
					}
					printchar(c, clr, &buf)
				}
				fmt.Fprintf(&buf, " ")
			}

			fmt.Fprintf(&buf, "\n")
		}
		fmt.Fprintf(&buf, "\n")
		fmt.Print(buf.String())
	}
}

func numNodes(g *Game, depth int) int {
	if depth <= 0 {
		return 1
	}

	mvs := g.LegalMoves()

	// i := len(mvs)

	i := 0

	for _, m := range mvs {
		gC := g.Copy()
		applyMoveUnchecked(&gC, m)
		i += numNodes(&gC, depth-1)
	}

	return i
}

func NumNodes(fen string, depth int) (int, error) {
	g, err := ParseFEN(fen)

	if err != nil {
		return 0, err
	}

	return numNodes(g, depth), nil
}
