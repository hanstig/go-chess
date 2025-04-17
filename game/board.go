package game

type PieceType string
type Color string

const (
	Pawn   PieceType = "P"
	Rook   PieceType = "R"
	Knight PieceType = "N"
	Bishop PieceType = "B"
	King   PieceType = "K"
	Queen  PieceType = "Q"

	White Color = "w"
	Black Color = "b"
)

type Piece struct {
	Type  PieceType
	Color Color
}

type Square struct {
	X int
	Y int
}

func inBounds(s Square) bool {
	return s.X >= 0 && s.X < 8 && s.Y >= 0 && s.Y < 8
}

type Board [8][8]*Piece

func copyBoard(b *Board) Board {
	var newBoard Board
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

func NewBoard() Board {
	return Board{}
}
