package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func (g *game) ToFen() string {
	fen := ""

	for y := range 8 {
		blanks := 0
		for x := range 8 {
			if g.board[y][x] == nil {
				blanks++
			} else {
				if blanks != 0 {
					fen += strconv.Itoa(blanks)
				}

				c := string(g.board[y][x].pieceType)
				if g.board[y][x].color == black {
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

	fen += string(g.turn)

	fen += " "

	cr := g.castlingRights

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

	if g.enPassantTarget == nil {
		fen += "-"
	} else {
		s := *g.enPassantTarget
		ept := ""
		ept += string(byte('a' + s.x))

		ept += strconv.Itoa(7 - s.y + 1)

		fen += ept
	}

	fen += " "

	fen += strconv.Itoa(g.halfMoveClock)
	fen += " "
	fen += strconv.Itoa(g.fullMoveCount)

	return fen
}

func ParseFEN(fen string) (*game, error) {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid FEN: expected 6 fields")
	}

	game := &game{}
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
				game.board[y][xpos] = p
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
		game.turn = white
	case "b":
		game.turn = black
	default:
		return nil, fmt.Errorf("invalid turn: %s", fields[1])
	}

	// Castling rights
	if strings.Contains(fields[2], "K") {
		game.castlingRights.whiteK = true
	}
	if strings.Contains(fields[2], "Q") {
		game.castlingRights.whiteQ = true
	}
	if strings.Contains(fields[2], "k") {
		game.castlingRights.blackK = true
	}
	if strings.Contains(fields[2], "q") {
		game.castlingRights.blackQ = true
	}

	// En passant
	if fields[3] != "-" {
		x := int(fields[3][0] - 'a')
		y := 7 - int(fields[3][1]-'1')
		game.enPassantTarget = &square{y: y, x: x}
	}

	// Halfmove clock
	hmc, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("Invalid halfmove clock: %v", err)
	}
	game.halfMoveClock = hmc

	// Fullmove number
	fmc, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("Invalid fullmove number: %v", err)
	}
	game.fullMoveCount = fmc

	return game, nil
}

func fenCharToPiece(c rune) (*piece, error) {
	color := white
	if unicode.IsLower(c) {
		color = black
		c = unicode.ToUpper(c)
	}

	var pieceType pieceType
	switch c {
	case 'P':
		pieceType = pawn
	case 'N':
		pieceType = knight
	case 'B':
		pieceType = bishop
	case 'R':
		pieceType = rook
	case 'Q':
		pieceType = queen
	case 'K':
		pieceType = king
	default:
		return nil, fmt.Errorf("invalid piece character: %c", c)
	}

	return &piece{pieceType: pieceType, color: color}, nil
}
