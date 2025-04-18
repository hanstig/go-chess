package chess

type pieceType string
type color string

const (
	pawn   pieceType = "P"
	rook   pieceType = "R"
	knight pieceType = "N"
	bishop pieceType = "B"
	king   pieceType = "K"
	queen  pieceType = "Q"

	white color = "w"
	black color = "b"
)

type piece struct {
	pieceType pieceType
	color     color
}

type square struct {
	x int
	y int
}

func (s square) inBounds() bool {
	return s.x >= 0 && s.x < 8 && s.y >= 0 && s.y < 8
}

type board [8][8]*piece

func copyBoard(b *board) board {
	var newBoard board
	for y := range 8 {
		for x := range 8 {
			if b[y][x] != nil {
				p := *b[y][x]
				newBoard[y][x] = &p
			}
		}
	}
	return newBoard
}
