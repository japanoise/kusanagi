package main

import "fmt"

var Vector [8][8]int = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},               // empty
	{0, 0, 0, 0, 0, 0, 0, 0},               // pawn - handled specially.
	{+21, +12, -8, -19, -21, -12, +8, +19}, // N
	{+11, -9, -11, +9, 0, 0, 0, 0},         // B
	{+10, +1, -10, -1, 0, 0, 0, 0},         // R
	{+10, +11, +1, -9, -10, -11, -1, +9},   // Q
	{+10, +11, +1, -9, -10, -11, -1, +9},   // K
	{0, 0, 0, 0, 0, 0, 0, 0},               // ?
}

var Slide [8]bool = [8]bool{false, false, false, true, true, true, false, false}

const (
	MoveQuiet byte = iota
	MoveDoublePush
	MoveCapture
)

type Move struct {
	From    byte
	To      byte
	Kind    byte
	Promote byte
	Score   int16
}

func pawnmove(b *Board, i int, retval []Move) []Move {
	var PawnPush, DoublePush int
	CanDouble := false
	if b.ToMove == BLACK {
		PawnPush = i - 10
		DoublePush = i - 20
		CanDouble = i/10 == 8
	} else {
		PawnPush = i + 10
		DoublePush = i + 20
		CanDouble = i/10 == 3
	}
	if GetPiece(b.Data[PawnPush]) == EMPTY {
		retval = append(retval, Move{byte(i),
			byte(PawnPush), MoveQuiet, EMPTY, 0})
		if CanDouble && GetPiece(b.Data[DoublePush]) ==
			EMPTY {
			retval = append(retval, Move{byte(i),
				byte(DoublePush), MoveDoublePush, EMPTY, 0})
		}
	}
	retval = pawncap(b, i, retval, PawnPush-1)
	retval = pawncap(b, i, retval, PawnPush+1)
	return retval
}

func pawncap(b *Board, i int, retval []Move, place int) []Move {
	if OnBoard(place) && GetPiece(b.Data[place]) != EMPTY &&
		GetSide(b.Data[place]) != b.ToMove {
		retval = append(retval, Move{byte(i),
			byte(place), MoveCapture, EMPTY, 0})
	}
	return retval
}

func squareattacked(b *Board, i int) bool {
	var PawnPush int
	if b.ToMove == BLACK {
		PawnPush = i + 10
	} else {
		PawnPush = i - 10
	}
	if (GetSide(b.Data[PawnPush-1]) != b.ToMove && GetPiece(b.Data[PawnPush-1]) == PAWN) || (GetSide(b.Data[PawnPush+1]) != b.ToMove && GetPiece(b.Data[PawnPush+1]) == PAWN) {
		return true
	}
	for dir := 0; dir < 8; dir++ {
		if Vector[QUEEN][dir] == 0 {
			break
		}
		from := i
		for {
			to := from + Vector[QUEEN][dir]
			piece := GetPiece(b.Data[to])
			if b.Data[to] == OFFBOARD || (piece != EMPTY && GetSide(b.Data[to]) == b.ToMove) {
				break
			} else if GetPiece(b.Data[to]) == QUEEN {
				return true
			} else if piece == ROOK && Vector[QUEEN][dir] == 10 || Vector[QUEEN][dir] == -10 || Vector[QUEEN][dir] == 1 || Vector[QUEEN][dir] == -1 {
				return true
			} else if piece == BISHOP && Vector[QUEEN][dir] == 1 || Vector[QUEEN][dir] == -1 || Vector[QUEEN][dir] == 9 || Vector[QUEEN][dir] == -9 {
				return true
			}
			from = to
		}
	}
	return false
}

func quietmove(b *Board, i int, retval []Move) []Move {
	piece := GetPiece(b.Data[i])
	for dir := 0; dir < 8; dir++ {
		if Vector[piece][dir] == 0 {
			break
		}
		from := i
		for {
			to := from + Vector[piece][dir]
			if b.Data[to] != OFFBOARD {
				if GetPiece(b.Data[to]) == EMPTY {
					retval = append(retval, Move{byte(i),
						byte(to), MoveQuiet, EMPTY, 0})
					if Slide[piece] {
						from = to
					} else {
						break
					}
				} else if GetSide(b.Data[to]) != b.ToMove {
					retval = append(retval, Move{byte(i),
						byte(to), MoveCapture, EMPTY, 0})
					break
				} else {
					break
				}
			} else {
				break
			}
		}
	}
	return retval
}

func MoveGen(b *Board) []Move {
	retval := make([]Move, 0, 32)
	for i := A1; i <= H8; i++ {
		if !OnBoard(i) || GetPiece(b.Data[i]) == EMPTY || GetSide(b.Data[i]) != b.ToMove {
			continue
		}
		if GetPiece(b.Data[i]) == PAWN {
			retval = pawnmove(b, i, retval)
		} else {
			retval = quietmove(b, i, retval)
		}
	}
	return retval
}

func MakeMove(b *Board, m *Move) {
	b.EnPassant = INVALID
	b.Data[m.To] = b.Data[m.From]
	b.Data[m.From] = EMPTY
	switch m.Kind {
	case MoveQuiet:
		/* Do nothing */
	case MoveDoublePush:
		if b.ToMove == BLACK {
			b.EnPassant = int(m.From - 10)
		} else {
			b.EnPassant = int(m.From + 10)
		}
	}
	b.ToMove ^= BLACK
}

func (m Move) String() string {
	return fmt.Sprint("{From: ", IndexToAlgebraic(int(m.From)), " to: ",
		IndexToAlgebraic(int(m.To)), " type: ", m.Kind, "}")
}

func DoPerft(depth int) uint64 {
	board, _ := Parse(START)
	return Perft(depth, board)
}

func Perft(depth int, board *Board) uint64 {
	if depth == 0 {
		return 1
	}
	var nodes uint64 = 0
	moves := MoveGen(board)
	for _, move := range moves {
		boardc := *board
		MakeMove(&boardc, &move)
		nodes += Perft(depth-1, &boardc)
	}
	return nodes
}
