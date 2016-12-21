package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

/*
WARNING: HACK APPROACHING!

The suggested way to store squares is with bit fiddling. So a square is a
byte, with the following structure:

000XCPPP

Where 0 is junk, X is validity, C is colour, and PPP is piece data.
*/

const (
	EMPTY byte = iota
	PAWN
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

const (
	WHITE byte = 0x00
	BLACK byte = 0x08
)

const (
	ONBOARD  byte = 0x00
	OFFBOARD byte = 0x10
)

type Board struct {
	/* Mailbox style, 10x12 board. */
	Data   [120]byte
	ToMove byte
}

const (
	A1 int = 21
	H8 int = 98
)

const START string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func ClearBoard(b *Board) {
	for i := 0; i < 120; i++ {
		if i < A1 || i > H8 {
			b.Data[i] = OFFBOARD
		} else if i%10 == 0 || i%10 == 9 {
			b.Data[i] = OFFBOARD
		} else {
			b.Data[i] = ONBOARD
		}
	}
}

func Parse(fen string) (*Board, error) {
	b := new(Board)
	ClearBoard(b)
	rank := 7
	file := 0
	stage := 0
	for _, runeValue := range fen {
		switch stage {
		case 0:
			/* Fill the data */
			if runeValue >= '1' && runeValue <= '8' {
				inc, _ := strconv.Atoi(string(runeValue))
				file += inc
			} else if runeValue == '/' {
				rank -= 1
				file = 0
			} else if runeValue == ' ' {
				stage++
			} else {
				switch unicode.ToUpper(runeValue) {
				case 'P':
					b.Data[CartesianToIndex(file, rank)] |= PAWN
				case 'N':
					b.Data[CartesianToIndex(file, rank)] |= KNIGHT
				case 'B':
					b.Data[CartesianToIndex(file, rank)] |= BISHOP
				case 'R':
					b.Data[CartesianToIndex(file, rank)] |= ROOK
				case 'Q':
					b.Data[CartesianToIndex(file, rank)] |= QUEEN
				case 'K':
					b.Data[CartesianToIndex(file, rank)] |= KING
				default:
					return nil, errors.New("Unexpected character in board data")
				}
				if unicode.IsLower(runeValue) {
					b.Data[CartesianToIndex(file, rank)] |= BLACK
				}
				file += 1
			}
		case 1:
			/* Get who's to play next */
			switch runeValue {
			case 'w':
				b.ToMove = WHITE
			case 'b':
				b.ToMove = BLACK
			case ' ':
				stage++
			default:
				return nil, errors.New("Unexpected character for active colour")
			}
		}
	}
	return b, nil
}

func PrintBoard(b *Board) string {
	retval := ""
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			retval +=
				ByteToString(b.Data[CartesianToIndex(file, rank)])
		}
		retval += "\n"
	}
	return retval
}

func CartesianToIndex(file, rank int) int {
	return 21 + (10 * rank) + file
}

func ByteToString(b byte) string {
	if ByteIsOffboard(b) {
		return ""
	}
	retval := ""
	switch GetPiece(b) {
	case EMPTY:
		return "."
	case PAWN:
		retval = "P"
	case KNIGHT:
		retval = "N"
	case BISHOP:
		retval = "B"
	case ROOK:
		retval = "R"
	case QUEEN:
		retval = "Q"
	case KING:
		retval = "K"
	default:
		return "?"
	}
	if IsBlack(b) {
		return strings.ToLower(retval)
	}
	return retval
}

func ByteIsOffboard(b byte) bool {
	return b&OFFBOARD == OFFBOARD
}

func GetPiece(b byte) byte {
	return b & 0x07
}

func IsBlack(b byte) bool {
	return b&BLACK == BLACK
}
