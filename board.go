package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	FORCE byte = 0x0F
)

const (
	ONBOARD  byte = 0x00
	OFFBOARD byte = 0x10
)

const (
	CASTLEWK byte = 0x01
	CASTLEWQ byte = 0x02
	CASTLEBK byte = 0x04
	CASTLEBQ byte = 0x08
)

var (
	Clock      time.Duration
	TimeRepeat int
	TimePerTC  time.Duration
	TimeInc    time.Duration
)

type Board struct {
	/* Mailbox style, 10x12 board. */
	Data      [120]byte
	ToMove    byte
	Castle    byte
	EnPassant byte
	WhiteKing byte
	BlackKing byte
	Moves     int
	PieceList [32]byte
}

const (
	A1 byte = 21
	H8 byte = 98
)

const START string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const INVALID byte = 0

var CASTLEMASK [120]byte

func InitCastleMask() {
	for sq := 0; sq < 120; sq++ {
		CASTLEMASK[sq] = 15 // all castle rights const?
	}

	sq, _ := AlgebraicToIndex("a1")

	CASTLEMASK[sq] ^= CASTLEWQ

	sq, _ = AlgebraicToIndex("h1")

	CASTLEMASK[sq] ^= CASTLEWK

	sq, _ = AlgebraicToIndex("e1")

	CASTLEMASK[sq] ^= CASTLEWK | CASTLEWQ

	sq, _ = AlgebraicToIndex("a8")

	CASTLEMASK[sq] ^= CASTLEBQ

	sq, _ = AlgebraicToIndex("h8")

	CASTLEMASK[sq] ^= CASTLEBK

	sq, _ = AlgebraicToIndex("e8")

	CASTLEMASK[sq] ^= CASTLEBK | CASTLEBQ
}

func ClearBoard(b *Board) {
	InitCastleMask()
	b.EnPassant = 1
	b.Castle = 0
	var i byte
	for i = 0; i < 120; i++ {
		if OnBoard(i) {
			b.Data[i] = ONBOARD
		} else {
			b.Data[i] = OFFBOARD
		}
	}
	for i = 0; i < 32; i++ {
		b.PieceList[i] = OFFBOARD
	}
}

func fenfillboard(b *Board, field string) error {
	var rank byte = 7
	var file byte = 0
	var idx byte = 0
	for _, runeValue := range field {
		/* Fill the data */
		if runeValue >= '1' && runeValue <= '8' {
			inc, _ := strconv.Atoi(string(runeValue))
			file += byte(inc)
		} else if runeValue == '/' {
			rank -= 1
			file = 0
		} else {
			sq := CartesianToIndex(file, rank)
			switch unicode.ToUpper(runeValue) {
			case 'P':
				b.Data[sq] |= PAWN
				b.PieceList[idx] = sq
			case 'N':
				b.Data[sq] |= KNIGHT
				b.PieceList[idx] = sq
			case 'B':
				b.Data[sq] |= BISHOP
				b.PieceList[idx] = sq
			case 'R':
				b.Data[sq] |= ROOK
				b.PieceList[idx] = sq
			case 'Q':
				b.Data[sq] |= QUEEN
				b.PieceList[idx] = sq
			case 'K':
				b.Data[sq] |= KING
				b.PieceList[idx] = sq
			default:
				return errors.New("Unexpected character in board data")
			}
			if unicode.IsLower(runeValue) {
				b.Data[sq] |= BLACK
			}
			file += 1
			idx += 1
		}
	}
	return nil
}

func fenwho(b *Board, field string) error {
	if field == "w" {
		b.ToMove = WHITE
	} else if field == "b" {
		b.ToMove = BLACK
	} else {
		return errors.New("Unexpected character for active colour")
	}
	return nil
}

func fencastling(b *Board, field string) error {
	for _, runeValue := range field {
		/* Castling */
		switch runeValue {
		case '-':
			b.Castle = 0
		case 'K':
			b.Castle |= CASTLEWK
		case 'Q':
			b.Castle |= CASTLEWQ
		case 'k':
			b.Castle |= CASTLEBK
		case 'q':
			b.Castle |= CASTLEBQ
		default:
			return errors.New("Unexpected character for castling")
		}
	}
	return nil
}

func fenenpassant(b *Board, field string) {
	epindex, err := AlgebraicToIndex(field)
	if err == nil {
		b.EnPassant = epindex
	} else {
		b.EnPassant = INVALID
	}
}

func Parse(fen string) (*Board, error) {
	b := new(Board)
	ClearBoard(b)
	fields := strings.Split(fen, " ")
	if len(fields) < 4 {
		return nil, errors.New("missing fields from fen string")
	}
	err := fenfillboard(b, fields[0])
	if err != nil {
		return nil, err
	}
	err = fenwho(b, fields[1])
	if err != nil {
		return nil, err
	}
	err = fencastling(b, fields[2])
	if err != nil {
		return nil, err
	}
	fenenpassant(b, fields[3])
	b.WhiteKing, _ = FindKing(b, WHITE)
	b.BlackKing, _ = FindKing(b, BLACK)
	return b, nil
}

func PrintBoard(b *Board) string {
	retval := ""
	var rank, file byte
	for rank = 7; rank != 255; rank-- {
		for file = 0; file < 8; file++ {
			retval +=
				ByteToString(b.Data[CartesianToIndex(file, rank)])
		}
		retval += "\n"
	}
	return retval
}

func CartesianToIndex(file, rank byte) byte {
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

func GetSide(b byte) byte {
	return b & BLACK
}

func IsBlack(b byte) bool {
	return GetSide(b) == BLACK
}

func OnBoard(i byte) bool {
	return i >= A1 && i <= H8 && !(i%10 == 0 || i%10 == 9)
}

func IndexToCartesian(index byte) (byte, byte) {
	var file, rank byte
	file = (index % 10) - 1
	rank = (index / 10) - 2
	return file, rank
}

func CartesianToAlgebraic(file, rank byte) string {
	afile := rune(byte('a') + file)
	arank := rank + 1
	return fmt.Sprintf("%c%d", afile, arank)
}

func IndexToAlgebraic(i byte) string {
	return CartesianToAlgebraic(IndexToCartesian(i))
}

func FindKing(b *Board, colour byte) (byte, error) {
	for king := A1; king <= H8; king++ {
		if b.Data[king] != OFFBOARD && GetPiece(b.Data[king]) == KING &&
			GetSide(b.Data[king]) == colour {
			return king, nil
		}
	}
	return INVALID, errors.New("Couldn't find the king")
}

func Illegal(b *Board) bool {
	king, err := GetKing(b, b.ToMove^BLACK)
	if err != nil {
		return true
	}
	return squareattacked(b, king, b.ToMove)
}

func AlgebraicToCartesian(a string) (byte, byte, error) {
	var rank, file byte
	if len(a) != 2 {
		return 0, 0, errors.New(fmt.Sprint("algebraic string", a, "was too long!"))
	}
	for i, runeValue := range a {
		if i == 1 {
			if runeValue >= '1' && runeValue <= '8' {
				irank, _ := strconv.Atoi(string(runeValue))
				rank = byte(irank)
			} else {
				return 0, 0, errors.New(fmt.Sprint("algebraic string", a, "has invalid rank!"))
			}
		} else {
			if runeValue >= 'a' && runeValue <= 'h' {
				file = byte(runeValue - 'a')
			} else {
				return 0, 0, errors.New(fmt.Sprint("algebraic string", a, "has invalid file!"))
			}
		}
	}
	return file, rank - 1, nil
}

func AlgebraicToIndex(a string) (byte, error) {
	file, rank, err := AlgebraicToCartesian(a)
	if err != nil {
		return INVALID, err
	}
	return CartesianToIndex(file, rank), nil
}

func GetKing(b *Board, side byte) (byte, error) {
	var retval byte
	if side == BLACK {
		retval = b.BlackKing
	} else {
		retval = b.WhiteKing
	}
	if retval == INVALID {
		return INVALID, errors.New("No king on the board")
	} else {
		return retval, nil
	}
}

func CanCastle(b *Board, color, side byte) bool {
	var flag byte
	if color == BLACK {
		if side == QUEEN {
			flag = CASTLEBQ
		} else {
			flag = CASTLEBK
		}
	} else {
		if side == QUEEN {
			flag = CASTLEWQ
		} else {
			flag = CASTLEWK
		}
	}
	return b.Castle&flag != 0
}

func FindPiece(b *Board, target byte) (byte, error) {
	for idx, sq := range b.PieceList {
		if sq == target {
			return byte(idx), nil
		}
	}
	return OFFBOARD, errors.New("Square is not in piece list")
}
