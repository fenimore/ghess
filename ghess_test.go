package ghess

import (
	"fmt"
	"testing"
	//"github.com/polypmer/ghess"
)

func TestNewBoard(t *testing.T) {
	game := NewBoard()
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if game.Position() != expected {
		t.Error("Unexpected FEN:", game.Position())
	}

}

// This test takes a while (like almost 3 seconds).
func TestPgnLoad(t *testing.T) {
	hist := `
[Event "State Ch."]
[Site "New York, USA"]
[Date "1910.??.??"]
[Round "?"]
[White "Capablanca"]
[Black "Jaffe"]
[Result "1-0"]
[ECO "D46"]
[Opening "Queen's Gambit Dec."]
[Annotator "Reinfeld, Fred"]
[WhiteTitle "GM"]
[WhiteCountry "Cuba"]
[BlackCountry "United States"]

1. d4 d5 2. Nf3 Nf6 3. e3 c6 4. c4 e6 5. Nc3 Nbd7 6. Bd3 Bd6
7. O-O O-O 8. e4 dxe4 9. Nxe4 Nxe4 10. Bxe4 Nf6 11. Bc2 h6
12. b3 b6 13. Bb2 Bb7 14. Qd3 g6 15. Rae1 Nh5 16. Bc1 Kg7
17. Rxe6 Nf6 18. Ne5 c5 19. Bxh6+ Kxh6 20. Nxf7+ 1-0`
	// 6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42
	var err error
	// TODO: Something is werid here in Bc1 move wtf
	// FIXME
	game := NewBoard()
	err = game.LoadPgn(hist)
	if err != nil {
		t.Error("Unable to parse pgn")
	}
	flaw := `
1. d4 d5 2. Nf3 Nf6 3. e3 c6 4. c4 e6 5. Nc3 Nbd7 6. Bd3 Bd6
7. O-O O-O 8. e4 dxe4 9. Nxe4 Nxe4 10. Bxe4 Nf6 11. Bc2 h6
12. b3 b6 13. Bb2 Bb7 14. Qd3 g6 15. Rae1 Nh Bc1 Kg7
17. Rxe6 Nf6 18. Ne5 c5 19. Bxh6+ Kxh6 20. Nxf7+ 1-0`
	game = NewBoard()
	err = game.LoadPgn(flaw)
	if err == nil {
		t.Error("This PGN shouldn't be parsed")
	}
}

func TestStandard(t *testing.T) {
	game := NewBoard()
	o := "e2"
	d := "e4"
	var err error
	err = game.ParseStand(o, d)
	if err != nil {
		t.Error("This move doesn't make any sense.")
	}
	if game.Position() != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Error("FEN looks funny ")
	}
}

func TestTurnChange(t *testing.T) {
	game := NewBoard()
	//var err error
	_ = game.Move(24, 44)
	if game.toMove != "b" {
		t.Error("Turn did not change")
	}
	s, _ := MiniMaxPruning(0, 2, GetState(&game))
	_ = game.Move(s.Init[0], s.Init[1])
	if game.toMove != "w" {
		t.Error("Minimax does not change turn?")
	}
}

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

func TestLoadPgnAndCheck(t *testing.T) {
	hist := `
[White "Amor"]
[Black "Caput"]
1. Nf3 Nc6 2. d4 d5 3. c4 e6 4. e3 Nf6 5. Nc3 Be7 6. a3 O-O 7. b4 a6 8. Be2 Re8 9. O-O Bf8 10. c5 g6 11. b5 axb5 12. Bxb5 Bd7 13. h3 Na5 14. Bd3 Nc6 15. Rb1 Qc8 16. Nb5 e5 17. Be2 e4 18. Ne1 h6 19. Nc2 g5 20. f3 exf3 21. Bxf3 g4 22. hxg4 Bxg4 23. Nxc7 Qxc7 24. Bxg4 Nxg4 25. Qxg4+ Bg7 26. Nb4 Nxb4 27. Rxb4 Ra6 28. Rf5 Re4 29. Qh5 Rg6 30. Qh3 Qc8 31. Qf3 Qd7 32. Rb2 Bxd4 33. exd4 Re1+ 34. Kh2 Rxc1 35. Qxd5 Qe7 36. g3 Qc7 37. Rf4 b6 38. a4 Rg5 39. cxb6 Rxd5 40. bxc7 Rxc7 41. Rb5 Rc2+`
	// 6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42
	var err error
	game := NewBoard()
	err = game.LoadPgn(hist)
	if err != nil {
		t.Error("But that's a legal PGN")
	}
	info := game.Stats()
	ch := game.PlayerCheck()
	if !ch {
		t.Error("But that's should be check")
	}
	if info["move"] != "42" {
		t.Error("But that's should be move 42")
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
		fmt.Println(game.StringWhite())
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

func TestEvaluateZero(t *testing.T) {
	game := NewBoard()
	if game.Evaluate() != 0 {
		t.Error("Init Position should be egal")
	}
}

func TestPruning(t *testing.T) {
	game := NewBoard()
	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	state_prun, _ := MiniMaxPruning(0, 2, s)
	state_mini, _ := MiniMax(0, 2, s)
	if state_mini.Init[0] != state_prun.Init[0] {
		t.Error("Unexpected")
	}
	if state_mini.Init[1] != state_prun.Init[1] {
		t.Error("Unmatched")
	}

}

func TestDictionaryAttack(t *testing.T) {
	game := NewBoard()
	s := GetState(&game)
	nextState, _ := MiniMaxPruning(0, 2, s)
	if nextState.Init[0] != 24 && nextState.Init[1] != 44 {
		t.Error(nextState, "Dictionary ERror")
	}

}

func TestDraw(t *testing.T) {
	game := NewBoard()
	err := game.Move(12, 33)
	if err != nil {
		fmt.Println(err)
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

func TestCheckMateMiniMax(t *testing.T) {
	game := NewBoard()
	// One more to checkmate
	fen := `4k3/8/8/8/8/7r/6r1/1K6 b - - 0 2`
	err := game.LoadFen(fen)
	if err != nil {
		t.Error(err)
	}
	s := GetState(&game)
	nxt, err := MiniMaxPruning(0, 2, s)
	if err != nil {
		t.Error(err)
	}
	err = game.Move(nxt.Init[0], nxt.Init[1])
	if err != nil {
		t.Error(err)
	}
	if !game.Checkmate {
		t.Error("Not Checkmate, Come ON")
	}
}

/**********************************
Chess Problems!!!
***********************************/
func TestChessProblemsMateInTwo(t *testing.T) {
	game := NewBoard()
	// Mate in Two
	// Morphy Verse Duke of Brunswick
	// 1. Qb8+ Nxb8 2. Rd8#
	fen := `4kb1r/p2n1ppp/4q3/4p1B1/4P3/1Q6/PPP2PPP/2KR4 w k - 1 1`
	// Alekhine vs Fahardo
	//1. Rxf6+ Nxf6 2. g5#
	//fen := `4r3/pbpn2n1/1p1prp1k/8/2PP2PB/P5N1/2B2R1P/R5K1 w - - 0 1`
	//Monterinas vs Euwe
	//1... Be3+ 2. Qxe3 Qg4#
	//fen := `7r/p3ppk1/3p4/2p1P1Kp/2Pb4/3P1QPq/PP5P/R6R b - - 0 1`
	err := game.LoadFen(fen)
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(game.StringWhite())

	nxt, err := MiniMaxPruning(0, 3, GetState(&game))
	if err != nil {
		t.Error(err)
	}
	err = game.Move(nxt.Init[0], nxt.Init[1])
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(game.StringWhite())
	// White Response
	nxt, err = MiniMaxPruning(0, 3, GetState(&game))
	if err != nil {
		t.Error(err)
	}

	err = game.Move(nxt.Init[0], nxt.Init[1])
	if err != nil {
		t.Error(err)
	}
	///fmt.Println(game.StringWhite())
	// Mate?
	nxt, err = MiniMaxPruning(0, 3, GetState(&game))
	if err != nil {
		t.Error(err)
	}
	err = game.Move(nxt.Init[0], nxt.Init[1])
	if err != nil {
		t.Error(err)
	}
	if !game.Checkmate {
		t.Error("Engine didn't solve problem")
	}
}

// func TestChessProblemsMateInThree(t *testing.T) {
//	game := NewBoard()
//	// Mate in Three
//	// Needs 5 ply
//	// Alekhine vs Freeman
//	fen := `4Rnk1/pr3ppp/1p3q2/5NQ1/2p5/8/P4PPP/6K1 w - - 1 1`
//	// Solution:

//	err := game.LoadFen(fen)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(game.StringWhite())

//	nxt, err := MiniMaxPruning(0, 5, GetState(&game))
//	if err != nil {
//		t.Error(err)
//	}
//	err = game.Move(nxt.Init[0], nxt.Init[1])
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(game.StringWhite())
//	// White Response
//	nxt, err = MiniMaxPruning(0, 5, GetState(&game))
//	if err != nil {
//		t.Error(err)
//	}

//	err = game.Move(nxt.Init[0], nxt.Init[1])
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(game.StringWhite())
//	// Mate?
//	nxt, err = MiniMaxPruning(0, 5, GetState(&game))
//	if err != nil {
//		t.Error(err)
//	}
//	err = game.Move(nxt.Init[0], nxt.Init[1])
//	if err != nil {
//		t.Error(err)
//	}

//	fmt.Println(game.StringWhite())
//	// White Response
//	nxt, err = MiniMaxPruning(0, 5, GetState(&game))
//	if err != nil {
//		t.Error(err)
//	}

//	err = game.Move(nxt.Init[0], nxt.Init[1])
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(game.StringWhite())
//	// Mate?
//	nxt, err = MiniMaxPruning(0, 5, GetState(&game))
//	if err != nil {
//		t.Error(err)
//	}
//	err = game.Move(nxt.Init[0], nxt.Init[1])
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(game.StringWhite())

//	if !game.Checkmate {
//		t.Error("Engine didn't solve problem")
//	}
// }

/**********************************
Examples
***********************************/

func Example_NewBoard() {
	game := NewBoard()
	fmt.Print(game.Position())
	// Output:
	// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
}

func ExampleBoard_ParseStand() {
	game := NewBoard()
	o := "e2"
	d := "e4"
	var err error
	err = game.ParseStand(o, d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(game.Position())
	err = game.ParseStand(o, d)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseStand("a1", "h8")
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseStand("g8", "g6")
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseStand("g8", "f6")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(game.Position())
	// Output:
	// rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1
	// Empty square
	// Not your turn
	// Illegal Knight Move
	// rnbqkb1r/pppppppp/5n2/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2
}

func ExampleBoard_LoadFen() {
	game := NewBoard()
	fen := "6Q1/8/8/p7/k7/5p2/1K6/8 w ---- - 0 5"
	var err error
	err = game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}
	fen = "6Q1/8/8/p7/k7/5p2/1K6/A w ---- - 0 5"
	err = game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}
	fen = "6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42"
	err = game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}
	fmt.Println(game.Check)

	//Output:
	// Success
	// Invalid FEN
	// Success
	// true
}

func ExampleEmpassantAndDisambigPawn() {
	game := NewBoard()
	fen := "rnbqkbnr/pp1ppppp/8/4P3/2p5/8/PPPP1PPP/RNBQKBNR w KQkq - 0 3"
	_ = game.LoadFen(fen)
	fmt.Println("Empassant:")
	fmt.Println(game.empassant)
	fmt.Println(string(game.board[37]), string(game.board[47]))
	game.ParseMove("b4")
	fmt.Println(game.empassant)
	fmt.Println(string(game.board[37]), string(game.board[47]))
	game.ParseMove("cxb3")
	fmt.Println(game.empassant)
	fmt.Println(string(game.board[37]), string(game.board[47]))
	game.ParseMove("cxb3")
	fmt.Println("C column:")
	fmt.Println(string(game.board[26])) // c pawn should be empty

	// Output:
	// Empassant:
	// 0
	// . .
	// 47
	// . P
	// 0
	// p .
	// C column:
	// .

}

func ExamplePgnDisambigRook() {
	fen := "1nbqkbnr/1pppppp1/r7/p6p/P6P/7R/1PPPPPP1/RNBQKBN1 w -Qk- - 0 4"
	game := NewBoard()
	_ = game.LoadFen(fen)
	err := game.ParseMove("Rha3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(game.board[31]))

	// Output:
	// .
}

func ExamplePgnDisambigRookAttack() {
	hist := `6k1/8/5p2/1R1r1R2/P2P2Pp/7K/2r5/8 w - - 0 45`
	game := NewBoard()
	err := game.LoadFen(hist)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseMove("Rfxd5")
	fmt.Println(string(game.board[53]))

	hist = `6k1/8/3r1p2/5R2/P2P2Pp/7K/5r2/5R2 w ---- - 0 47`
	err = game.LoadFen(hist)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseMove("R1xf2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(game.board[13]))

	// Output:
	// .
	// .
}

func ExamplePgnDisambigKnight() {
	// Column Disambig
	fen := `2bqr1k1/r3bp1p/p1np1np1/1p2p3/3NP1PN/1B2BP2/PPPQ3P/2KR3R w - - 0 17`
	game := NewBoard()
	err := game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseMove("Nhf5")
	fmt.Println(string(game.board[14]))
	// Row Disambig
	fen = `r5nr/Np1k3p/n4Q2/8/1bB4q/N1P1p3/PP3PPP/R3K2R w KQ - 1 19`
	game = NewBoard()
	err = game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseMove("N7b5")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(game.board[83]))

	// Output:
	// .
	// .
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

func ExampleBoard_SearchValid() {
	game := NewBoard()
	o, d := game.SearchValid()
	fmt.Println(o)
	fmt.Println(d)
	fen := "6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42"
	_ = game.LoadFen(fen)
	o, d = game.SearchValid()
	fmt.Println(o)
	fmt.Println(d)

	// Output:
	//[12 12 17 17 21 21 22 22 23 23 24 24 25 25 26 26 27 27 28 28]
	//[33 31 38 36 31 41 32 42 33 43 34 44 35 45 36 46 37 47 38 48]
	//[21 21 21 43]
	//[31 11 12 23]

}

/**********************************
Benchmarks MiniMax
***********************************/

func BenchmarkMidGamePruningDepth2(b *testing.B) {
	// Very short time
	game := NewBoard()
	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 2, s)
		if err != nil {
			fmt.Println("Minimax 2 error: ", err)
		}
	}

}

func BenchmarkOpeningPruningDepth2(b *testing.B) {
	// Opening position doesn't count,
	// cause of the dictionary attack
	game := NewBoard()
	fen := "r1bqkbnr/ppp2ppp/2np4/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 2, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkOpeningPruningDepth3(b *testing.B) {
	// Opening position doesn't count,
	// cause of the dictionary attack
	// Seems to be about 14 seconds
	game := NewBoard()
	fen := "r1bqkbnr/ppp2ppp/2np4/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 3, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth3(b *testing.B) {
	// Seems to be about four seconds
	game := NewBoard()
	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 3, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth3v2(b *testing.B) {
	// Seems to be about four seconds
	game := NewBoard()
	fen := "rn1qkb1r/1p3ppp/p2pbn2/4p3/4P3/1NN1BP2/PPP3PP/R2QKB1R b KQkq - 0 7"
	_ = game.LoadFen(fen)
	//s :=
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 3, GetState(&game))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkOpeningPruningDepth4(b *testing.B) {
	// Opening position doesn't count,
	// cause of the dictionary attack
	// Seems to be about one and a half  minute for depth four
	game := NewBoard()
	fen := "r1bqkbnr/ppp2ppp/2np4/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 4, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth4(b *testing.B) {
	// Seems to be about 75 seconds
	game := NewBoard()
	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 4, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth4v2(b *testing.B) {
	// Seems to be about four seconds
	game := NewBoard()
	fen := "rn1qkb1r/1p3ppp/p2pbn2/4p3/4P3/1NN1BP2/PPP3PP/R2QKB1R b KQkq - 0 7"
	_ = game.LoadFen(fen)
	//s :=
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 4, GetState(&game))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Woah going to five
func BenchmarkOpeningPruningDepth5(b *testing.B) {
	// Opening position doesn't count,
	// cause of the dictionary attack

	game := NewBoard()
	fen := "r1bqkbnr/ppp2ppp/2np4/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 5, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth5(b *testing.B) {

	game := NewBoard()
	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
	_ = game.LoadFen(fen)
	s := GetState(&game)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 5, s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMidGamePruningDepth5v2(b *testing.B) {
	// Seems to be about four seconds
	game := NewBoard()
	fen := "rn1qkb1r/1p3ppp/p2pbn2/4p3/4P3/1NN1BP2/PPP3PP/R2QKB1R b KQkq - 0 7"
	_ = game.LoadFen(fen)
	//s :=
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 5, GetState(&game))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// // Woah going to five
// func BenchmarkOpeningPruningDepth6(b *testing.B) {
//	// Opening position doesn't count,
//	// cause of the dictionary attack

//	game := NewBoard()
//	fen := "r1bqkbnr/ppp2ppp/2np4/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4"
//	_ = game.LoadFen(fen)
//	s := GetState(&game)
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		_, err := MiniMaxPruning(0, 6, s)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
// }

// func BenchmarkMidGamePruningDepth6(b *testing.B) {

//	game := NewBoard()
//	fen := "r1bqkb1r/1p3ppp/p1n2n2/3p4/8/1N1B4/PPP2PPP/RNBQ1RK1 w kq - 0 9"
//	_ = game.LoadFen(fen)
//	s := GetState(&game)
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		_, err := MiniMaxPruning(0, 6, s)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
// }

/*
func BenchmarkMidGamePruningDepth6v2(b *testing.B) {
	// Seems to be about four seconds
	game := NewBoard()
	fen := "rn1qkb1r/1p3ppp/p2pbn2/4p3/4P3/1NN1BP2/PPP3PP/R2QKB1R b KQkq - 0 7"
	_ = game.LoadFen(fen)
	//s :=
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := MiniMaxPruning(0, 6, GetState(&game))
		if err != nil {
			fmt.Println(err)
		}
	}
}
*/
