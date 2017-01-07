package ghess

import (
	"fmt"
	"testing"
)

/* Piece validation check */

func TestQueenValid(t *testing.T) {
	game := NewBoard()
	var err error
	// weird Queen surrounded by pawns
	fen := `rnbqkbnr/8/8/2ppp3/2pQp3/2ppp3/P4PPP/R5KR w KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen Error")
	}
	err = game.Move(45, 27)
	if err == nil {
		t.Error("Queen can't move there!!")
	}
	err = game.Move(45, 67)
	if err == nil {
		t.Error("Queen can't move there!!")
	}
	err = game.Move(45, 36)
	if err != nil {
		t.Error("Queen should be able to move here")
	}
	fen = `rnbqkbnr/ppNppppp/8/8/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen Error")
	}
	err = game.Move(85, 76)
	if err != nil {
		t.Error("Queen should be able to move here")
	}
	fen = `rnbqkbnr/ppNppppp/8/8/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen Error")
	}

	fen = `rnbqkbnr/pp1ppppp/8/8/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}
	err = game.Move(85, 76)
	if err != nil {
		t.Error("Queen should be able to move here")
	}
	fen = `rnbqkbnr/pp1ppppp/8/8/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}
	err = game.Move(85, 67)
	if err != nil {
		t.Error("Queen should be able to move here")
	}

	fen = `r2q1kbr/ppNN1ppp/4p3/3p4/4P3/8/PPPP1PPP/R1BQKB1R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}
	err = game.Move(85, 75)
	if err != nil {
		t.Error("Queen should be able to move here")
	}

}

func TestBishopValid(t *testing.T) {
	game := NewBoard()
	var err error
	fen := `rnbqkbnr/pp1p1ppp/8/8/1P1p4/2B5/PP2PPPP/RN1QKBNR w KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen Error")
	}
	err = game.Move(36, 54)
	if err == nil {
		t.Error("Bishop can't move there!!")
	}
	err = game.Move(36, 58)
	if err == nil {
		t.Error("Bishop can't move there!!")
	}
	err = game.Move(36, 47)
	if err == nil {
		t.Error("Bishop can't move there!!")
	}
	err = game.Move(36, 45)
	if err != nil {
		t.Error("Bishop should move there!!")
	}
}

func TestPawnValid(t *testing.T) {
	fen := `rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 0 1`
	game := NewBoard()
	var err error
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen Error")
	}

	err = game.Move(23, 43)
	if err == nil {
		t.Error("Illegal Pawn move")
	}
	err = game.Move(23, 34)
	if err == nil {
		t.Error("There's no attack there..")
	}
	err = game.Move(24, 44)
	if err != nil {
		t.Error("But that's a legal move")
	}

}

func TestCannotCastleThroughCheck(t *testing.T) {
	game := NewBoard()
	hist := `r3k2r/ppp4p/n4Q1n/8/1bB3q1/N1P1p2N/PP3PPP/R3K2R w KQkq - 0 14`
	var err error
	err = game.LoadFen(hist)
	if err != nil {
		t.Error("But that's a legal FEN")
	}
	//_ = game.PlayerCheck()
	err = game.ParseMove("O-O-O")
	if err == nil {
		t.Error("Shouldn't be allowed to castle")
	}
}

func TestCannotCastleInCheck(t *testing.T) {
	game := NewBoard()
	fen := `rnb1kbnr/pppppppp/8/2q5/3P4/4P3/PPP2PPP/RNBQK2R b KQkq - 0 20`
	err := game.LoadFen(fen)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	err = game.Move(56, 47)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	// Now I should be in check
	// Unable to castle
	err = game.Move(14, 11)
	if err == nil {
		t.Error("No castling in Check")
	}
}

func TestCheckWithPawn(t *testing.T) {
	game := NewBoard()
	hist := "1. e4 e5 2. Ke2 Qf6 3. Kd3 Nh6 4. Kc4"
	var err error
	err = game.LoadPgn(hist)
	if err != nil {
		t.Error("Shouldn't be error")
	}
	err = game.ParseMove("b5")
	if err != nil {
		t.Error("Shouldn't be error")
	}
	ch := game.PlayerCheck()
	if !ch {
		t.Error("Should be check")
	}
}

func TestMoveIntoCheck(t *testing.T) {
	game := NewBoard()
	hist := `1. e4 e5 2. Qf3 Qg5 3. Qxf7 Ke7`
	var err error
	err = game.LoadPgn(hist)
	if err == nil {
		t.Error("Valid")
	}
	fen := `r1b1kbnr/p1qp1ppp/npp5/4p3/4P1K1/5P2/PPPP2PP/RNBQ1BNR w kq - 1 6`
	_ = game.LoadFen(fen)
	err = game.ParseMove("Kg5")
	if err != nil {
		t.Error("Hould be valid")
	}
	fen = `	rnbqk1nr/ppp2ppp/5b2/3pp3/8/1K1P4/PPP1PPPP/RNBQ1BNR w kq - 4 5`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Can't load fen")
	}
	err = game.ParseMove("Kb4")
	if err != nil {
		t.Error("Valid")
	}
	// Output:
	// Cannot move into Check
}

func TestDraw(t *testing.T) {
	game := NewBoard()
	err := game.Move(12, 33)
	if err != nil {
		t.Error(err)
	}
	_ = game.Move(82, 63)
	_ = game.Move(33, 12)
	_ = game.Move(63, 82)
	//fmt.Println(game.Score, game.Draw)
	_ = game.Move(12, 33)
	_ = game.Move(82, 63)
	_ = game.Move(33, 12)
	_ = game.Move(63, 82)
	//fmt.Println(game.Score, game.Draw)
	_ = game.Move(12, 33)
	_ = game.Move(82, 63)
	_ = game.Move(33, 12)
	_ = game.Move(63, 82)
	//fmt.Println(game.history)
	if !game.Draw {
		t.Error("Should be a draw")
	}
}

func ExampleCheckMate() {
	// Must call PlayerCheckMate
	game := NewBoard()
	fen := "6Q1/8/8/p7/k7/5p2/1K6/8 w ---- - 0 5"
	_ = game.LoadFen(fen)
	_ = game.PlayerCheckMate()
	fmt.Println(game.Checkmate)
	game.ParseMove("Qc4")
	_ = game.PlayerCheckMate()
	fmt.Println(game.Checkmate)

	// Output:
	// false
	// true
}
