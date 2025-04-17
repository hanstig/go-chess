package game

import "fmt"

type Move struct {
	From      Square
	To        Square
	Promotion *PieceType
}

func (m Move) deltas() (int, int) {
	return m.To.X - m.From.X, m.To.Y - m.From.Y
}

func ApplyMove(g *Game, m Move) error {
	if !IsLegalMove(g, m) {
		return fmt.Errorf("Illegal move!")
	}

	return applyMoveUnchecked(g, m)
}

func applyMoveUnchecked(g *Game, m Move) error {
	from := m.From
	to := m.To
	b := &(g.Board)

	piece := b[from.Y][from.X]
	// target := b[to.Y][to.X]
	if piece == nil {
		return fmt.Errorf("no piece at source square")
	}

	setEnPassantThisTurn := false
	// g.EnPassantTarget = nil
	if piece.Type == Pawn {
		if to.Y == 0 || to.Y == 7 {
			if m.Promotion == nil {
				return fmt.Errorf("Pawn moved to end without promotion")
			}

			prm := *m.Promotion

			if !(prm == Queen || prm == Rook || prm == Knight || prm == Bishop) {
				return fmt.Errorf("Pawn moved to end with invalid promotion: %v", prm)
			}
			piece.Type = prm
		}

		dx, dy := m.deltas()
		if abs(dy) == 2 {
			s := Square{
				X: m.From.X,
				Y: m.From.Y + dy/2,
			}
			g.EnPassantTarget = &s
			setEnPassantThisTurn = true
		}

		if g.EnPassantTarget != nil && to == *g.EnPassantTarget {
			b[from.Y][from.X+dx] = nil
		}
	} else if piece.Type == King {
		dx, _ := m.deltas()

		if abs(dx) == 2 {
			rookFromX := 0
			if dx > 0 {
				rookFromX = 7
			}

			rook := g.Board[from.Y][rookFromX]

			rookToX := from.X + (dx / 2)
			b[from.Y][rookToX] = rook
			b[from.Y][rookFromX] = nil
		}

		if piece.Color == White {
			g.CastlingRights.whiteK = false
			g.CastlingRights.whiteQ = false
		} else {
			g.CastlingRights.blackK = false
			g.CastlingRights.blackQ = false
		}

	} else if piece.Type == Rook {
		// Update castling rights
		if piece.Color == White {
			if from.X == 0 && from.Y == 7 {
				g.CastlingRights.whiteQ = false
			} else if from.X == 7 && from.Y == 7 {
				g.CastlingRights.whiteK = false
			}
		} else {
			if from.X == 0 && from.Y == 0 {
				g.CastlingRights.blackQ = false
			} else if from.X == 7 && from.Y == 0 {
				g.CastlingRights.blackK = false
			}
		}
	}

	if !setEnPassantThisTurn {
		g.EnPassantTarget = nil
	}

	if m.To == (Square{Y: 7, X: 7}) {
		g.CastlingRights.whiteK = false
	} else if m.To == (Square{Y: 7, X: 0}) {
		g.CastlingRights.whiteQ = false
	} else if m.To == (Square{Y: 0, X: 7}) {
		g.CastlingRights.blackK = false
	} else if m.To == (Square{Y: 0, X: 0}) {
		g.CastlingRights.blackQ = false
	}

	b[to.Y][to.X] = piece
	b[from.Y][from.X] = nil

	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}

	return nil
}
