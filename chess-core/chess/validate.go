package chess

import "fmt"

func (g *game) isPseudoLegalMove(m move) bool {
	from := m.from
	to := m.to
	if !from.inBounds() || !to.inBounds() {
		return false
	}

	piece := g.board[from.y][from.x]
	if piece == nil || piece.color != g.turn {
		return false
	}

	target := g.board[to.y][to.x]
	if target != nil && target.color == g.turn {
		return false
	}

	switch piece.pieceType {
	case pawn:
		return g.validatePawnMove(m)
	case rook:
		return g.validateRookMove(m)
	case knight:
		return g.validateKnightMove(m)
	case bishop:
		return g.validateBishopMove(m)
	case queen:
		return g.validateQueenMove(m)
	case king:
		return g.validateKingMove(m)
	}

	return false
}

func (g *game) validatePawnMove(m move) bool {
	prm := queen
	m.promotion = &prm

	dir := -1
	homeRow := 6
	if g.turn == black {
		dir = 1
		homeRow = 1
	}

	from := m.from
	to := m.to

	dx, dy := m.deltas()

	if dx == 0 {
		// single step
		if from.y+dir == to.y && g.board[to.y][to.x] == nil {
			return true
		}

		// double step
		if from.y == homeRow && from.y+2*dir == to.y && g.board[to.y][to.x] == nil {
			return isPathClear(g, m)
		}
	} else if abs(dy) == 1 && abs(dx) == 1 && from.y+dir == to.y {
		// en passant
		return g.board[to.y][to.x] != nil || (g.enPassantTarget != nil && to == *g.enPassantTarget)
	}

	return false
}

func (g *game) validateRookMove(m move) bool {
	return isStraightMove(m) && isPathClear(g, m)
}

func (g *game) validateKnightMove(m move) bool {
	dx, dy := m.deltas()
	return dx*dx+dy*dy == 5
}

func (g *game) validateBishopMove(m move) bool {
	return isDiagonalMove(m) && isPathClear(g, m)
}

func (g *game) validateKingMove(m move) bool {
	dx, dy := m.deltas()
	piece := g.board[m.from.y][m.from.x]

	// Castling
	if abs(dx) == 2 && dy == 0 {

		allowed := true
		var rookSquare square

		if piece.color == white {
			rookSquare.y = 7
			if dx < 0 {
				rookSquare.x = 0
				allowed = g.castlingRights.whiteQ
			} else {
				rookSquare.x = 7
				allowed = g.castlingRights.whiteK
			}
		} else { // piece.Color == Black
			rookSquare.y = 0
			if dx < 0 {
				rookSquare.x = 0
				allowed = g.castlingRights.blackQ
			} else {
				rookSquare.x = 7
				allowed = g.castlingRights.blackK
			}
		}

		if !allowed {
			return false
		}

		// can now assume king and rook on correct squares

		// can't castle out of check
		g.swapTurn()
		valid := !g.hasCheckMate()
		g.swapTurn()
		if !valid {
			return false
		}

		// Path between king and rook needs to be clear
		kingToRook := move{
			from: m.from,
			to:   rookSquare,
		}

		if !isPathClear(g, kingToRook) {
			return false
		}

		// cant castle through check
		gCopy := g.copy()
		kingToGap := move{
			from: m.from,
			to:   square{y: m.from.y, x: m.from.x + (dx / 2)},
		}

		gCopy.applyMoveUnchecked(kingToGap)

		if gCopy.hasCheckMate() {
			return false
		}

		return true
	}

	return dx*dx+dy*dy <= 2
}

func (g *game) validateQueenMove(m move) bool {
	return (isDiagonalMove(m) || isStraightMove(m)) && isPathClear(g, m)
}

func (g *game) IsLegalMove(m move) bool {
	// NOTE: Is there a speed difference here?
	// prm := queen
	// m.promotion = &prm

	m = setPromotion(m, queen)

	if !g.isPseudoLegalMove(m) {
		return false
	}

	newGame := g.copy()

	err := newGame.applyMoveUnchecked(m)

	if err != nil {
		return false
	}

	return !newGame.hasCheckMate()
}

func (g *game) findKing(c color) (*square, error) {
	for y := range 8 {
		for x := range 8 {
			p := g.board[y][x]
			if p != nil && p.pieceType == king && p.color == c {
				return &square{x: x, y: y}, nil
			}
		}
	}
	return nil, fmt.Errorf("No %v king with this color exists", c)
}

func (g *game) hasCheckMate() bool {
	enemyColor := white
	if g.turn == white {
		enemyColor = black
	}

	enemyKingSquare, err := g.findKing(enemyColor)

	if err != nil {
		// No enemy king, act like player already has checkmate
		return true
	}

	b := &(g.board)

	for y := range 8 {
		for x := range 8 {
			p := b[y][x]

			if p == nil || p.color != g.turn {
				continue
			}

			move := move{
				from: square{y: y, x: x},
				to:   *enemyKingSquare,
			}

			if g.isPseudoLegalMove(move) {
				return true
			}
		}
	}
	return false
}
