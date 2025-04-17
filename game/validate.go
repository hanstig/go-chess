package game


func IsPseudoLegalMove(g *Game, m Move) bool {
	from := m.From
	to := m.To
	if !inBounds(from) || !inBounds(to) {
		return false
	}

	piece := g.Board[from.Y][from.X]
	if piece == nil || piece.Color != g.Turn {
		return false
	}

	target := g.Board[to.Y][to.X]
	if target != nil && target.Color == g.Turn {
		return false
	}

	switch piece.Type {
	case Pawn:
		return validatePawnMove(g, m)
	case Rook:
		return validateRookMove(g, m)
	case Knight:
		return validateKnightMove(g, m)
	case Bishop:
		return validateBishopMove(g, m)
	case Queen:
		return validateQueenMove(g, m)
	case King:
		return validateKingMove(g, m)
	}

	return false
}

func validatePawnMove(g *Game, m Move) bool {
	prm := Queen
	m.Promotion = &prm

	dir := -1
	homeRow := 6
	if g.Turn == Black {
		dir = 1
		homeRow = 1
	}

	from := m.From
	to := m.To

	dx, dy := m.deltas()

	if dx == 0 {
		// single step
		if from.Y+dir == to.Y && g.Board[to.Y][to.X] == nil {
			return true
		}

		// double step
		if from.Y == homeRow && from.Y+2*dir == to.Y && g.Board[to.Y][to.X] == nil {
			return isPathClear(g, m)
		}
	}

	// diagonal capture or en passant
	if abs(dy) == 1 && abs(dx) == 1 && from.Y+dir == to.Y {
		return g.Board[to.Y][to.X] != nil || (g.EnPassantTarget != nil && to == *g.EnPassantTarget)
	}

	return false
}

func validateRookMove(g *Game, m Move) bool {
	return isStraightMove(m) && isPathClear(g, m)
}

func validateKnightMove(g *Game, m Move) bool {
	dx, dy := m.deltas()
	return dx*dx+dy*dy == 5
}

func validateBishopMove(g *Game, m Move) bool {
	return isDiagonalMove(m) && isPathClear(g, m)
}

func validateKingMove(g *Game, m Move) bool {
	dx, dy := m.deltas()
	piece := g.Board[m.From.Y][m.From.X]

	// Castling
	if abs(dx) == 2 && dy == 0 {

		allowed := true
		var rookSquare Square

		if piece.Color == White {
			rookSquare.Y = 7
			if dx < 0 {
				rookSquare.X = 0
				allowed = g.CastlingRights.whiteQ
			} else {
				rookSquare.X = 7
				allowed = g.CastlingRights.whiteK
			}
		} else { // piece.Color == Black
			rookSquare.Y = 0
			if dx < 0 {
				rookSquare.X = 0
				allowed = g.CastlingRights.blackQ
			} else {
				rookSquare.X = 7
				allowed = g.CastlingRights.blackK
			}
		}

		if !allowed {
			return false
		}

		// can now assume king and rook on correct squares

		// cant castle out of check
		gCopy := g.Copy()
		if gCopy.Turn == White {
			gCopy.Turn = Black
		} else {
			gCopy.Turn = White
		}
		if hasCheckMate(&gCopy) {
			return false
		}

		// Path between king and rook needs to be clear
		kingToRook := Move{
			From: m.From,
			To:   rookSquare,
		}
		if !isPathClear(g, kingToRook) {
			return false
		}

		// cant castle through check
		gCopy = g.Copy()
		kingToGap := Move{
			From: m.From,
			To:   Square{Y: m.From.Y, X: m.From.X + (dx / 2)},
		}
		applyMoveUnchecked(&gCopy, kingToGap)
		if hasCheckMate(&gCopy) {
			return false
		}

		return true
	}

	return dx*dx+dy*dy <= 2
}

func validateQueenMove(g *Game, m Move) bool {
	return (isDiagonalMove(m) || isStraightMove(m)) && isPathClear(g, m)
}

func IsLegalMove(g *Game, m Move) bool {

	prm := Queen
	m.Promotion = &prm

	if !IsPseudoLegalMove(g, m) {
		return false
	}

	// if m.From.X == 0 && m.From.Y == 1 && m.To.X == 0 && m.To.Y == 0 {
	// 	fmt.Println("is pseudo legal")
	// 	PrintMoves(g,[]Move{m})
	// }

	newGame := g.Copy()

	err := applyMoveUnchecked(&newGame, m)

	if err != nil {
		return false
	}

	return !hasCheckMate(&newGame)
}

func findKing(g *Game, c Color) *Square {
	for y := range 8 {
		for x := range 8 {
			p := g.Board[y][x]
			if p != nil && p.Type == King && p.Color == c {
				return &Square{X: x, Y: y}
			}
		}
	}
	return nil
}

func hasCheckMate(g *Game) bool {
	enemyColor := White
	if g.Turn == White {
		enemyColor = Black
	}

	enemyKingSquare := findKing(g, enemyColor)

	b := &(g.Board)

	for y := range 8 {
		for x := range 8 {
			p := b[y][x]

			if p == nil || p.Color != g.Turn {
				continue
			}

			move := Move{
				From: Square{Y: y, X: x},
				To:   *enemyKingSquare,
			}

			if IsPseudoLegalMove(g, move) {
				return true
			}
		}
	}
	return false
}
