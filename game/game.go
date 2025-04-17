package game

import (
	"fmt"
	"slices"
	"strings"
)

const (
	DefaultFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type Game struct {
	Board Board
	Turn  Color

	CastlingRights struct {
		whiteK bool
		whiteQ bool
		blackK bool
		blackQ bool
	}

	EnPassantTarget *Square
	HalfMoveClock   int
	FullMoveCount   int
}

func (g *Game) Copy() Game {
	newGame := *g

	newGame.Board = copyBoard(&(g.Board))

	if g.EnPassantTarget != nil {
		ept := *(g.EnPassantTarget)
		newGame.EnPassantTarget = &ept
	}

	return newGame
}

func NewGame() *Game {
	g, _ := ParseFEN(DefaultFen)
	return g
}

func (g *Game) HasLegalMoves() bool {
	b := &(g.Board)
	for fy := range 8 {
		for fx := range 8 {
			p := b[fy][fx]
			if p == nil || p.Color != g.Turn {
				continue
			}

			from := Square{Y: fy, X: fx}

			for ty := range 8 {
				for tx := range 8 {
					to := Square{X: tx, Y: ty}
					move := Move{From: from, To: to}

					if IsLegalMove(g, move) {
						return true
					}
				}
			}

		}
	}

	return false
}

func (g *Game) LegalMovesFrom(from Square) []Move {

	p := g.Board[from.Y][from.X]

	legalMoves := make([]Move, 0)

	if p == nil || p.Color != g.Turn {
		return legalMoves
	}

	for ty := range 8 {
		for tx := range 8 {
			to := Square{X: tx, Y: ty}
			move := Move{From: from, To: to}

			legal := IsLegalMove(g, move)

			if legal && (p.Type == Pawn && (ty == 0 || ty == 7)) {
				prm1 := Queen
				prm2 := Bishop
				prm3 := Rook
				prm4 := Knight

				m1 := move
				m1.Promotion = &prm1
				m2 := move
				m2.Promotion = &prm2
				m3 := move
				m3.Promotion = &prm3
				m4 := move
				m4.Promotion = &prm4

				legalMoves = slices.Concat(legalMoves, []Move{m1, m2, m3, m4})

			} else if legal {
				legalMoves = append(legalMoves, move)
			}
		}
	}

	return legalMoves
}

func (g *Game) LegalMoves() []Move {
	legalMoves := make([]Move, 0)
	b := &(g.Board)
	for fy := range 8 {
		for fx := range 8 {
			p := b[fy][fx]
			if p == nil || p.Color != g.Turn {
				continue
			}

			from := Square{Y: fy, X: fx}
			legalMoves = slices.Concat(legalMoves, g.LegalMovesFrom(from))
		}
	}
	return legalMoves
}

func (g *Game) Print() {
	for y := range 8 {
		for x := range 8 {
			c := "."
			if p := g.Board[y][x]; p != nil {
				c = string(p.Type)
				if p.Color == Black {
					c = strings.ToLower(string(c))
				}
			}
			fmt.Printf("%s ", c)
		}
		fmt.Println()
	}
	fmt.Println()
}
