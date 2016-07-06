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
- Queen Move validation
- FEN import
- Move history/ Undo
  - Save game history to board (not automatic)?
  - Save as two coordinates, with piece specifier

### Notes...

    rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
    rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1
    rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2
    rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2

