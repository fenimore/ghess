# go-chess || Ghess
A Golang chess engine, fumbling along...

    |Move:  3     Turn: b
    |Check: false Castle: KQkq
    |Mate:  false Score: *
    8: |♖||♘||♗||♕||♔|| ||♘||♖|
    7: |♙||♙||♙||♙|| ||♙||♙||♙|
    6: |░|| ||░|| ||░|| ||░|| |
    5: | ||░||♗||░||♙||░|| ||░|
    4: |░|| ||░||♟||♟|| ||░|| |
    3: | ||░|| ||░|| ||♞|| ||░|
    2: |♟||♟||♟|| ||░||♟||♟||♟|
    1: |♜||♞||♝||♛||♚||♝|| ||♜|
       :a::b::c::d::e::f::g::h:
    Black to move: 


# Instructions
`ghess.go` is a go package responsible for representing a chess board, parsing PGN input, and validating moves. In the `ui/` directory, `clichess.go` is a simple interface for testing.

After putting the source in `$GOPATH/src/github.com/polypmer/ghess/`, try

    go run ui/clichess.go

# Features
- Most rules are implemented:
  * No queen promotion
  * No queen PGN parse disambiguation
- PGN import export parse
- FEN import export parse
- Cli interface
- Random game!

## Search and Evaluate

- looks for all possible moves
- looks for all valid moves

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

    Coordinates:
    8: |88||87||86||85||84||83||82||81|
    7: |78||77||76||75||74||73||72||71|
    6: |68||67||66||65||64||63||62||61|
    5: |58||57||56||55||54||53||52||51|
    4: |48||47||46||45||44||43||42||41|
    3: |38||37||36||35||34||33||32||31|
    2: |28||27||26||25||24||23||22||21|
    1: |18||17||16||15||14||13||12||11|
       :a ::b ::c ::d ::e ::f ::g ::h :

- 120 bytes
- 11 - 18 1st rank
- 81 - 88 8th rank
- Does it seem backwards?

## TODO

### Basic rules

- Minor pawn promotion.

### Basic Functionality

- Variables should be exported, capitalized
- Checkmate should update PGN headers/history
- Flip board (UI problem)

### Search & Evaluate

- Look for all valid moves
   * There is a bug for castling possibilities
   * and likely empassant...
- Give all moves a score (mdr)

### Extra features

- Move history/ Undo
  * Save game history to board (not automatic)?
  * Save as two coordinates, with piece specifier

### Bugs:

- Surely, something?

### Notes...

    rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
    rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2
	8/5p2/8/p5P1/k7/8/1K6/8 w - - 0 1
	6Q1/8/8/p7/k7/5p2/1K6/8 w ---- - 0 5


### License;

    Fenimore Love - 2016 | GPL
