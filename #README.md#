# go-chess || Ghess
A Golang chess engine, fumbling along...

    Move: 3 | Castle: KQkq | Turn:  Black
    8: |♖||♘||♗||♕||♔||♗||·||♖|
    7: |♙||♙||♙||·||♙||♙||♙||♙|
    6: |·||·||·||·||·||♘||·||·|
    5: |·||·||·||♙||·||·||·||·|
    4: |·||·||♟||♟||·||·||·||·|
    3: |·||·||·||·||·||♞||·||·|
    2: |♟||♟||·||·||♟||♟||♟||♟|
    1: |♜||♞||♝||♛||♚||♝||·||♜|
       :a::b::c::d::e::f::g::h:
    Black to move: 

# Instructions
- ghess.go is the chess program. Run clichess.go, in the `ui/` directory in order to play a game of chess.


# Features
- Basic rules are implemented
  * Except:
  * Only queen promotion
  * no disambigious pgn parsing :/
- PGN import export parse
- FEN import export parse
- Cli interface
- random game!

## Search and Evaluate

- look for all possible moves
- look for all valid moves

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


## TODO

- More tests

### Basic rules

- Minor Pawn promotion

### Basic Functionality

- PGN parse errors
- FEN turn signification???
- Variables should be exported, capitalized
- Checkmate should update PGN headers
- Checkmate should # pgn

### Search & Evaluate

- Look for all valid moves
   * There is a bug for castling possibilities
   * and likely empassant...
- Give all moves a score

### Extra features

- Move history/ Undo
- Save game history to board (not automatic)?
- Save as two coordinates, with piece specifier

### Bugs:

- qe4 crashes
- Enfait, crashes when corrupt basic input...
- Load Fen doesn't print corrent turn (UI)
- Load FEN doesn't change check/checkmate values?


### Notes...

    rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1
    rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2
    rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2
	8/5p2/8/p5P1/k7/8/1K6/8 w - - 0 1
	6Q1/8/8/p7/k7/5p2/1K6/8 w ---- - 0 5

from the godoc docs
 
    Notice this comment is a complete sentence that begins with the name of the element it describes. This important convention allows us to generate documentation in a variety of formats, from plain text to HTML to UNIX man pages, and makes it read better when tools truncate it for brevity, such as when they extract the first line or sentence. 


### License;

Fenimore Love - 2016
GPL
