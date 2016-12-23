package main

import (
	"testing"
)

func TestMakeMovePawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(0, 2))
	from := byte(CartesianToIndex(0, 1))
	move := Move{from, to, MoveQuiet, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) != PAWN {
		t.Fail()
	}
}

func TestMakeMoveWhitePawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(0, 3))
	from := byte(CartesianToIndex(0, 1))
	move := Move{from, to, MoveDoublePush, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || board.EnPassant != CartesianToIndex(0, 2) {
		t.Fail()
	}
}

func TestMakeMoveBlackPawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	to := byte(CartesianToIndex(0, 4))
	from := byte(CartesianToIndex(0, 6))
	move := Move{from, to, MoveDoublePush, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || board.EnPassant != CartesianToIndex(0, 5) {
		t.Fail()
	}
}

func TestMakeMoveWhitePawnCapture(t *testing.T) {
	board, err := Parse("8/8/8/3p4/4P3/8/8/8 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(3, 4))
	from := byte(CartesianToIndex(4, 3))
	move := Move{from, to, MoveCapture, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || GetSide(board.Data[to]) != WHITE {
		t.Fail()
	}
}

func TestMakeMoveBlackPawnCapture(t *testing.T) {
	board, err := Parse("8/8/8/3p4/4P3/8/8/8 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(4, 3))
	from := byte(CartesianToIndex(3, 4))
	move := Move{from, to, MoveCapture, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || GetSide(board.Data[to]) != BLACK {
		t.Fail()
	}
}

func IsMoveInMoveList(t *testing.T, list []Move, from, to, kind byte) bool {
	for _, m := range list {
		if m.To == to && m.From == from && m.Kind == kind {
			return true
		}
	}
	return false
}

func TestMoveGenWhitePawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 1)), byte(CartesianToIndex(0, 2)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 6)), byte(CartesianToIndex(0, 5)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenWhitePawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 1)), byte(CartesianToIndex(0, 3)), MoveDoublePush) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 6)), byte(CartesianToIndex(0, 4)), MoveDoublePush) {
		t.Fail()
	}
}

func TestMoveGenSlider(t *testing.T) {
	board, err := Parse("4k3/8/8/1r6/8/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(1, 4)), byte(CartesianToIndex(1, 7)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenNonSlider(t *testing.T) {
	board, err := Parse("4k3/8/8/1r6/8/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 7)), byte(CartesianToIndex(4, 6)), MoveQuiet) {
		t.Fail()
	}
	if IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 7)), byte(CartesianToIndex(4, 3)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenWhitePawnCap(t *testing.T) {
	board, err := Parse("4k3/8/8/3p4/4P3/8/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 3)), byte(CartesianToIndex(3, 4)), MoveCapture) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnCap(t *testing.T) {
	board, err := Parse("4k3/8/8/3p4/4P3/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(3, 4)), byte(CartesianToIndex(4, 3)), MoveCapture) {
		t.Fail()
	}
}
