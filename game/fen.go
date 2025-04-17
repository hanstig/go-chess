package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func (g *Game) ToFen() string {

	fen := ""

	for y := range 8 {
		blanks := 0
		for x := range 8 {
			if g.Board[y][x] == nil {
				blanks++
			} else {
				if blanks != 0 {
					fen += strconv.Itoa(blanks)
				}

				c := string(g.Board[y][x].Type)
				if g.Board[y][x].Color == Black {
					c = strings.ToLower(c)
				}
				fen += c

				blanks = 0
			}
		}

		if blanks != 0 {
			fen += strconv.Itoa(blanks)
		}
		if y != 7 {
			fen += "/"
		}
	}

	fen += " "

	fen += string(g.Turn)

	fen += " "

	cr := g.CastlingRights

	if !cr.whiteK && !cr.whiteQ && !cr.blackK && !cr.blackQ {
		fen += "-"
	} else {
		if cr.whiteK {
			fen += "K"
		}
		if cr.whiteQ {
			fen += "Q"
		}
		if cr.blackK {
			fen += "k"
		}
		if cr.blackQ {
			fen += "q"
		}
	}

	fen += " "

	if g.EnPassantTarget == nil {
		fen += "-"
	} else {
		s := *g.EnPassantTarget
		ept := ""
		ept += string(byte('a' + s.X))
		
		ept += strconv.Itoa(7 - s.Y + 1)

		fen += ept
	}

	fen += " "

	fen += strconv.Itoa(g.HalfMoveClock)
	fen += " "
	fen += strconv.Itoa(g.FullMoveCount)

	return fen
}

func ParseFEN(fen string) (*Game, error) {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid FEN: expected 6 fields")
	}

	game := &Game{}
	boardRows := strings.Split(fields[0], "/")
	if len(boardRows) != 8 {
		return nil, fmt.Errorf("invalid board layout in FEN")
	}

	// TODO: Make this safe, can panic on some bad fens i think
	for y, row := range boardRows {
		xpos := 0
		for _, c := range row {
			p, err := fenCharToPiece(rune(c))

			if err == nil {
				game.Board[y][xpos] = p
				xpos += 1
				continue
			}

			num, err := strconv.Atoi(string(c))

			if err != nil || num <= 0 || num > 8 || xpos+num > 8 {
				return nil, errors.New("Invalid fen string, first field is bad")
			}

			xpos += num
		}
	}

	// Parse active color
	switch fields[1] {
	case "w":
		game.Turn = White
	case "b":
		game.Turn = Black
	default:
		return nil, fmt.Errorf("invalid turn: %s", fields[1])
	}

	// Castling rights
	if strings.Contains(fields[2], "K") {
		game.CastlingRights.whiteK = true
	}
	if strings.Contains(fields[2], "Q") {
		game.CastlingRights.whiteQ = true
	}
	if strings.Contains(fields[2], "k") {
		game.CastlingRights.blackK = true
	}
	if strings.Contains(fields[2], "q") {
		game.CastlingRights.blackQ = true
	}

	// En passant
	if fields[3] != "-" {
		x := int(fields[3][0] - 'a')
		y := 7 - int(fields[3][1]-'1')
		game.EnPassantTarget = &Square{Y: y, X: x}
	}

	// Halfmove clock
	hmc, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid halfmove clock: %v", err)
	}
	game.HalfMoveClock = hmc

	// Fullmove number
	fmc, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("invalid fullmove number: %v", err)
	}
	game.FullMoveCount = fmc

	return game, nil
}

func fenCharToPiece(c rune) (*Piece, error) {
	color := White
	if unicode.IsLower(c) {
		color = Black
		c = unicode.ToUpper(c)
	}

	var pieceType PieceType
	switch c {
	case 'P':
		pieceType = Pawn
	case 'N':
		pieceType = Knight
	case 'B':
		pieceType = Bishop
	case 'R':
		pieceType = Rook
	case 'Q':
		pieceType = Queen
	case 'K':
		pieceType = King
	default:
		return nil, fmt.Errorf("invalid piece character: %c", c)
	}

	return &Piece{Type: pieceType, Color: color}, nil
}
