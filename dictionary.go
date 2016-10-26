package ghess

var dict = make(map[string][2]int)

func dictionary() {
	/* e4 e5 */
	//1 e4
	dict["rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w"] = [2]int{24, 44}
	//1 e4 e5
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b"] = [2]int{74, 54}
	//2. Bc4
	dict["rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w"] = [2]int{13, 46}
	//2. Bc4 Nf6
	dict["rnbqkbnr/pppp1ppp/8/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR b"] = [2]int{82, 63}
	// 3. Nc3
	dict["rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR w"] = [2]int{17, 36}
	// 3. Nc3 Nc6
	dict["rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/2N5/PPPP1PPP/R1BQK1NR b"] = [2]int{87, 66}

	// alternative white move 2. Nf3
	// 2 Nf3 Nc6
	dict["rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b"] = [2]int{87, 66}
	// 3 Bc4
	dict["r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R b"] = [2]int{83, 56}

	/* d4 d5 */
	//1 d4 d5
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR b"] = [2]int{75, 55}
	//2 Nf3
	dict["rnbqkbnr/pppppppp/8/8/4P3/8/PPP1PPPP/RNBQKBNR w"] = [2]int{12, 33}

	/*  Nf3 */
	//1 Nf3 d5
	dict["rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b"] = [2]int{75, 55}
	// 2. d4
	dict["rnbqkbnr/ppp1pppp/8/3p4/8/5N2/PPPPPPPP/RNBQKB1R w"] = [2]int{25, 45}

	/* Sicilian


	1   e4 c5
	2.   Nf3 d6
	3.   d4 cxd4
	4.   Nxd4 Nf6
	5.   Nc3 g6
	6.   Be3 Bg7
	7.   f3 O-O
	*/
	// after c5 -> Nf3
	dict["rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w"] = [2]int{12, 33}
	// 2. Nf3 -> d6
	dict["rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b"] = [2]int{75, 65}
	// 3. -> d4
	dict["rnbqkbnr/pp2pppp/3p4/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R w"] = [2]int{25, 45}
	// 3. d4 -> cxd4
	dict["rnbqkbnr/pp2pppp/3p4/2p5/3PP3/5N2/PPP2PPP/RNBQKB1R b"] = [2]int{56, 45}
	// 4. -> Nxd4
	dict["rnbqkbnr/pp2pppp/3p4/8/3pP3/5N2/PPP2PPP/RNBQKB1R w"] = [2]int{33, 45}
	// 4. Nxd4 -> g6
	dict["rnbqkbnr/pp2pppp/3p4/8/3NP3/8/PPP2PPP/RNBQKB1R b"] = [2]int{72, 62}
}
