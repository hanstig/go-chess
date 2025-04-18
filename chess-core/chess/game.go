package chess

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	defaultFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type game struct {
	board board
	turn  color

	castlingRights struct {
		whiteK bool
		whiteQ bool
		blackK bool
		blackQ bool
	}

	enPassantTarget *square
	halfMoveClock   int
	fullMoveCount   int
}

func (g *game) copy() game {
	newGame := *g

	newGame.board = copyBoard(&(g.board))

	if g.enPassantTarget != nil {
		ept := *(g.enPassantTarget)
		newGame.enPassantTarget = &ept
	}

	return newGame
}

func (g *game) hasLegalMoves() bool {
	b := &(g.board)
	for fy := range 8 {
		for fx := range 8 {
			p := b[fy][fx]
			if p == nil || p.color != g.turn {
				continue
			}

			from := square{y: fy, x: fx}

			for ty := range 8 {
				for tx := range 8 {
					to := square{x: tx, y: ty}
					move := move{from: from, to: to}

					if g.IsLegalMove(move) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (g *game) legalMovesFrom(from square) []move {
	p := g.board[from.y][from.x]

	legalMoves := make([]move, 0)

	if p == nil || p.color != g.turn {
		return legalMoves
	}

	for ty := range 8 {
		for tx := range 8 {
			to := square{x: tx, y: ty}
			m := move{from: from, to: to}
			legal := g.IsLegalMove(m)
			if legal && (p.pieceType == pawn && (ty == 0 || ty == 7)) {
				m1 := setPromotion(m, queen)
				m2 := setPromotion(m, bishop)
				m3 := setPromotion(m, rook)
				m4 := setPromotion(m, knight)
				legalMoves = slices.Concat(legalMoves, []move{m1, m2, m3, m4})
			} else if legal {
				legalMoves = append(legalMoves, m)
			}
		}
	}

	return legalMoves
}

func moveToString(m move) string {
	fx := string(byte('a' + m.from.x))
	fy := strconv.Itoa(7 + 1 - m.from.y)
	tx := string(byte('a' + m.to.x))
	ty := strconv.Itoa(7 + 1 - m.to.y)

	p := ""
	if m.promotion != nil {
		p = strings.ToLower(string(*m.promotion))
	}
	return fx + fy + tx + ty + p
}

func stringToSquare(str string) (square, error) {
	var s square
	s.x = int(str[0] - 'a')
	s.y = 7 - int(str[1]-'1')

	if !s.inBounds() {
		return square{}, fmt.Errorf("Invalid square string %v", str)
	}

	return s, nil
}

func stringToMove(str string) (move, error) {
	if len(str) != 4 && len(str) != 5 {
		return move{}, fmt.Errorf("Invalid move string %v", str)
	}

	from, err := stringToSquare(str[0:2])

	if err != nil {
		return move{}, fmt.Errorf("Invalid move string: %v", str)
	}

	to, err := stringToSquare(str[2:4])

	if err != nil {
		return move{}, fmt.Errorf("Invalid move string: %v", str)
	}

	m := move{from: from, to: to}
	if len(str) == 5 {
		pt := pieceType(strings.ToUpper(string(str[4])))
		if pt != queen && pt != rook && pt != bishop && pt != knight {
			return move{}, fmt.Errorf("Invalid move string: %v", str)
		}
		m = setPromotion(m, pt)
	}

	return m, nil
}

func (g *game) legalMoves() []move {
	legalMoves := make([]move, 0)
	b := &(g.board)
	for fy := range 8 {
		for fx := range 8 {
			p := b[fy][fx]
			if p == nil || p.color != g.turn {
				continue
			}

			from := square{y: fy, x: fx}
			legalMoves = slices.Concat(legalMoves, g.legalMovesFrom(from))
		}
	}
	return legalMoves
}

func (g *game) swapTurn() {
	if g.turn == black {
		g.turn = white
	} else {
		g.turn = black
	}
}
