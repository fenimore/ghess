package ghess

var dict = make(map[string][2]int)

func dictionary() {
	/* e4 e5 */
	//1 e4
	dict["rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"] = [2]int{24, 44}
	//1 e4 e5
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"] = [2]int{74, 54}
	//2. Bc4
	dict["rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"] = [2]int{13, 46}
	//2. Bc4 Nf6
	dict["rnbqkbnr/pppp1ppp/8/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR b KQkq - 0 2"] = [2]int{82, 63}
	// alternative white move 2. Nf3
	// 2 Nf3 Nc6
	dict["rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 2"] = [2]int{87, 66}
	// 3 Bc4
	dict["r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 0 3"] = [2]int{83, 56}

	/* d4 d5 */
	//1 d4 d5
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1"] = [2]int{75, 55}
	//2 Nf3
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR w KQkq d6 0 2"] = [2]int{12, 33}

	/*  Nf3 */
	//1 Nf3 d5
	dict["rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 0 1"] = [2]int{75, 55}
	// 2. d4
	dict["rnbqkbnr/ppp1pppp/8/3p4/8/5N2/PPPPPPPP/RNBQKB1R w KQkq d6 0 2"] = [2]int{25, 45}
}
