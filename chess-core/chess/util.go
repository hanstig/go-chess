package chess

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

func isStraightMove(m move) bool {
	dx, dy := m.deltas()
	return (dx == 0) != (dy == 0)
}

func isDiagonalMove(m move) bool {
	dx, dy := m.deltas()
	return abs(dx) == abs(dy) && dx != 0
}


func setPromotion(m move, p pieceType) move {
	m.promotion = &p
	return m
}

func (m move) deltas() (int, int) {
	return m.to.x - m.from.x, m.to.y - m.from.y
}

func isPathClear(g *game, m move) bool {
	dx, dy := m.deltas()
	dx, dy = sign(dx), sign(dy)

	x, y := m.from.x+dx, m.from.y+dy

	for x != m.to.x || y != m.to.y {
		if g.board[y][x] != nil {
			return false
		}

		x += dx
		y += dy
	}

	return true
}

