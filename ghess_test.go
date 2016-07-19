package ghess

import (
	"fmt"
	"testing"
	"strconv"
	//"github.com/polypmer/ghess"
)

func TestNewBoard(t *testing.T) {

}


func ExampleStartPosition() {
	game := NewBoard()
	fmt.Print(game.Position())
	// Output:
	// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
}

func ExampleLoadPgnAndCheck() {
	hist := `1. Nf3 Nc6 2. d4 d5 3. c4 e6 4. e3 Nf6 5. Nc3 Be7 6. a3 O-O 7. b4 a6 8. Be2 Re8 9. O-O Bf8 10. c5 g6 11. b5 axb5 12. Bxb5 Bd7 13. h3 Na5 14. Bd3 Nc6 15. Rb1 Qc8 16. Nb5 e5 17. Be2 e4 18. Ne1 h6 19. Nc2 g5 20. f3 exf3 21. Bxf3 g4 22. hxg4 Bxg4 23. Nxc7 Qxc7 24. Bxg4 Nxg4 25. Qxg4+ Bg7 26. Nb4 Nxb4 27. Rxb4 Ra6 28. Rf5 Re4 29. Qh5 Rg6 30. Qh3 Qc8 31. Qf3 Qd7 32. Rb2 Bxd4 33. exd4 Re1+ 34. Kh2 Rxc1 35. Qxd5 Qe7 36. g3 Qc7 37. Rf4 b6 38. a4 Rg5 39. cxb6 Rxd5 40. bxc7 Rxc7 41. Rb5 Rc2+`
	// 6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42
	var err error
	game := NewBoard()
	game, err = game.LoadPgn(hist)
	if err != nil {
		fmt.Println(err)
	}
	info := game.Stats()
	ch,_:= strconv.ParseBool(info["check"])
	if ch {
		fmt.Println("****Check!****")
	}
	fmt.Println(info["move"])
	// Output:
	// ****Check!****
	// 42
	
}

func ExampleCannotCastleThroughCheck() {
	game := NewBoard()
	hist := `r3k2r/ppp4p/n4Q1n/8/1bB3q1/N1P1p2N/PP3PPP/R3K2R w KQkq - 0 14`
	var err error
	err = game.LoadFen(hist)
	if err != nil {
		fmt.Println(err)
	}
	err = game.ParseMove("O-O-O")
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// Cannot Castle through check
}

func ExampleMoveIntoCheck() {
	game := NewBoard()
	hist := `1. e4 e5 2. Qf3 Qg5 3. Qxf7 Ke7`
	var err error
	game, err = game.LoadPgn(hist)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// Cannot move into Check
}

func ExampleLoadFen() {
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

	//Output:
	// Success
	// Invalid FEN
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
	fmt.Println(string(game.board[26]))// c pawn should be empty

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
	fmt.Println(string(game.board[83]))

	// Output:
	// .
	// .
}

func ExampleCheckMate() {
	game := NewBoard()
	fen := "6Q1/8/8/p7/k7/5p2/1K6/8 w ---- - 0 5"
	_ = game.LoadFen(fen)
	fmt.Println(game.checkmate)
	game.ParseMove("Qc4")
	fmt.Println(game.checkmate)


	// Output:
	// false
	// true
}

func ExampleSearchForValid() {
	game := NewBoard()
	fen := "6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42"
	_ = game.LoadFen(fen)
	fmt.Print(game.String())

		

}
