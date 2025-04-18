package chess

import (
	"fmt"
	"strings"
)

func NewGame() *game {
	g, err := ParseFEN(defaultFen)

	// Should never happen
	if err != nil {
		panic(err.Error())
	}

	return g
}

func NewGameFromFen(fen string) (*game, error) {
	g, err := ParseFEN(fen)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *game) ApplyMove(mv string) error {
	move, err := stringToMove(mv)

	if err != nil {
		return err
	}

	return g.applyMove(move)
}

func (g *game) LegalMoves() []string {
	mvs := g.legalMoves()
	mvStrings := make([]string, 0, len(mvs))
	for _, m := range mvs {
		mvStrings = append(mvStrings, moveToString(m))
	}
	return mvStrings
}

func (g *game) LegalMovesFrom(from string) ([]string, error) {
	f, err := stringToSquare(from)

	if err != nil {
		return nil, fmt.Errorf("Invalid From string: %v", from)
	}

	mvs := g.legalMovesFrom(f)

	mvStrings := make([]string, 0, len(mvs))
	for _, m := range mvs {
		mvStrings = append(mvStrings, moveToString(m))
	}

	return mvStrings, nil
}

func (g *game) GameOver() bool {
	return !g.hasLegalMoves()
}

func (g *game) Print() {
	b := &g.board

	for y := range 8 {
		for x := range 8 {
			p := b[y][x]
			if p == nil {
				fmt.Print(".")
			} else {
				c := string(p.pieceType)
				if p.color == black {
					c = strings.ToLower(c)
				}
				fmt.Print(c)
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
