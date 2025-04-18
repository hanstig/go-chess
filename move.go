package chess

import "fmt"

type move struct {
	from      square
	to        square
	promotion *pieceType
}

func (g *game) applyMove(m move) error {
	if !g.IsLegalMove(m) {
		return fmt.Errorf("Illegal move!")
	}

	return g.applyMoveUnchecked(m)
}

func (g *game) applyMoveUnchecked(m move) error {
	from := m.from
	to := m.to
	b := &(g.board)

	piece := b[from.y][from.x]
	if piece == nil {
		return fmt.Errorf("no piece at source square")
	}

	setEnPassantThisTurn := false
	if piece.pieceType == pawn {

		dx, dy := m.deltas()
		if to.y == 0 || to.y == 7 {
			// pawn promotion
			if m.promotion == nil {
				return fmt.Errorf("Pawn moved to end without promotion")
			}
			prm := *m.promotion
			if !(prm == queen || prm == rook || prm == knight || prm == bishop) {
				return fmt.Errorf("Pawn moved to end with invalid promotion: %v", prm)
			}
			piece.pieceType = prm
		} else if abs(dy) == 2 {
			// double step, set en passant
			s := square{
				x: m.from.x,
				y: m.from.y + dy/2,
			}
			g.enPassantTarget = &s
			setEnPassantThisTurn = true
		} else if g.enPassantTarget != nil && to == *g.enPassantTarget {
			b[from.y][from.x+dx] = nil
		}
	} else if piece.pieceType == king {
		dx, _ := m.deltas()

		if abs(dx) == 2 {
			rookFromX := 0
			if dx > 0 {
				rookFromX = 7
			}

			rook := g.board[from.y][rookFromX]

			rookToX := from.x + (dx / 2)
			b[from.y][rookToX] = rook
			b[from.y][rookFromX] = nil
		}

		if piece.color == white {
			g.castlingRights.whiteK = false
			g.castlingRights.whiteQ = false
		} else {
			g.castlingRights.blackK = false
			g.castlingRights.blackQ = false
		}

	} else if piece.pieceType == rook {
		// Update castling rights
		if piece.color == white {
			if from.x == 0 && from.y == 7 {
				g.castlingRights.whiteQ = false
			} else if from.x == 7 && from.y == 7 {
				g.castlingRights.whiteK = false
			}
		} else {
			if from.x == 0 && from.y == 0 {
				g.castlingRights.blackQ = false
			} else if from.x == 7 && from.y == 0 {
				g.castlingRights.blackK = false
			}
		}
	}

	if !setEnPassantThisTurn {
		g.enPassantTarget = nil
	}

	if m.to == (square{y: 7, x: 7}) {
		g.castlingRights.whiteK = false
	} else if m.to == (square{y: 7, x: 0}) {
		g.castlingRights.whiteQ = false
	} else if m.to == (square{y: 0, x: 7}) {
		g.castlingRights.blackK = false
	} else if m.to == (square{y: 0, x: 0}) {
		g.castlingRights.blackQ = false
	}

	b[to.y][to.x] = piece
	b[from.y][from.x] = nil

	g.swapTurn()

	return nil
}
