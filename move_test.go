package main

import (
	"testing"
)

func TestMoveGenWhitePawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if len(moves) == 0 {
		t.FailNow()
	}
	if moves[0].From != 31 && moves[0].To != 41 {
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
	if len(moves) == 0 {
		t.FailNow()
	}
	if moves[0].From != 81 && moves[0].To != 71 {
		t.Fail()
	}
}
