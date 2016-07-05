# go-chess || Ghess
Golang chess engine, fumbling along...

## Board struct
- A board object
- Game variables
- castle possibilities
- score
- to move
- move count
- fen string
- pgn string


## Bitmap
- 120 bytes
- 11 - 18 1st rank
- 81 - 88 8th rank

## TODO
- Check/checkmate is going to be a pain...
- King and Queen Move validation
- FEN import
- Store PGN
- Castle PGN input
- Print pgn history on crash
- And many more
- empassant
- Move history/ Undo
- Save game history to board (not automatic)

### Fen String examples:

rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1

rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1

rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2

rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2

### Print rotate??

	game := b.board
	p := b.pieces
	fmt.Print(p[string(game[89])], p[string(game[81])], 
		p[string(game[82])], p[string(game[83])], 
		p[string(game[84])], p[string(game[85])],
		p[string(game[86])], p[string(game[87])], 
		p[string(game[88])], "\n")
	fmt.Print(p[string(game[79])], p[string(game[71])], 
		p[string(game[72])],
		p[string(game[73])], p[string(game[74])],
		p[string(game[75])],
		p[string(game[76])], p[string(game[77])],
		p[string(game[78])], "\n")
	fmt.Print(p[string(game[69])], p[string(game[61])],
		p[string(game[62])],
		p[string(game[63])], p[string(game[64])],
		p[string(game[65])],
		p[string(game[66])], p[string(game[67])],
		p[string(game[68])], "\n")
	fmt.Print(p[string(game[59])], p[string(game[51])],
		p[string(game[52])],
		p[string(game[53])], p[string(game[54])],
		p[string(game[55])],p[string(game[56])],
		p[string(game[57])],
		p[string(game[58])], "\n")
	fmt.Print(p[string(game[49])], p[string(game[41])],
		p[string(game[42])],
		p[string(game[43])], p[string(game[44])],
		p[string(game[45])],
		p[string(game[46])], p[string(game[47])],
		p[string(game[48])], "\n")
	fmt.Print(p[string(game[39])], p[string(game[31])],
		p[string(game[32])],
		p[string(game[33])], p[string(game[34])],
		p[string(game[35])],
		p[string(game[36])], p[string(game[37])],
		p[string(game[38])], "\n")
	fmt.Print(p[string(game[29])], p[string(game[21])],
		p[string(game[22])],
		p[string(game[23])], p[string(game[24])],
		p[string(game[25])],
		p[string(game[26])], p[string(game[27])],
		p[string(game[28])], "\n")
	fmt.Print(p[string(game[19])], p[string(game[11])],
		p[string(game[12])],
		p[string(game[13])], p[string(game[14])],
		p[string(game[15])],
		p[string(game[16])], p[string(game[17])],
		p[string(game[18])], "\n")
	fmt.Print(string(game[91]), string(game[92]), string(game[93]),
		string(game[94]), string(game[95]), string(game[96]),
		string(game[97]), string(game[98]),
		string(game[99]), "\n")