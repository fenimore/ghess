# go-chess || Ghess
Golang chess engine, fumbling along...

## Board struct
- A board object
- Game variables
  * castle possibilities
  * score
  * to move
  * check
  * empassant
- Board display
  * pgnMap
  * pieceMap
  * pieces
- move count
- fen string
- pgn string
- pgn Headers

## Bitmap
- 120 bytes
- 11 - 18 1st rank
- 81 - 88 8th rank

## TODO
- Check/checkmate is going to be a pain...
- King and Queen Move validation
- FEN import
- Castle PGN input
- Move history/ Undo
- Save game history to board (not automatic)?

### Fen String examples:

rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1

rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1

rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2

rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2

### Print rotate??
// The pgnMap is broken if it is rotated. I don't understand this. It is likely very simple to understand. In any case, gonna keep the default String() method as white oriented, and so the pgnMap will be 18-11 not 11-18

This is the old pgnMap:
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"] = 11, 12, 13, 14, 15, 16, 17, 18
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"] = 21, 22, 23, 24, 25, 26, 27, 28
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"] = 31, 32, 33, 34, 35, 36, 37, 38
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"] = 41, 42, 43, 44, 45, 46, 47, 48
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"] = 51, 52, 53, 54, 55, 56, 57, 58
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"] = 61, 62, 63, 64, 65, 66, 67, 68
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"] = 71, 72, 73, 74, 75, 76, 77, 78
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"] = 81, 82, 83, 84, 85, 86, 87, 88
