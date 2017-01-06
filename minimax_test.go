package ghess

import (
	"fmt"
	"testing"
)

func TestDictionaryAttack(t *testing.T) {
	game := NewBoard()
	s := GetState(&game)
	nextState, _ := MiniMaxPruning(0, 2, s)
	if nextState.Init[0] != 24 && nextState.Init[1] != 44 {
		t.Error(nextState, "Dictionary ERror")
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
