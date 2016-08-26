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

- looks for all possible and valid moves


<hr>

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

<hr>

## TODO

1. more tests
2. Benchmarks for Search and Validation?
   * is it taking a while, or is it just me?

### Basic rules

- Minor pawn promotion.

### Basic Functionality

- Variables should be exported, capitalized
- Checkmate should update PGN headers/history
- Flip board (UI problem)

### Search & Evaluate

- Look for all valid moves
   * There is a bug for castling possibilities
- Give all moves a score (mdr)
   * Am I even capable of this?
   
### Extra features

- Move history/ Undo
  * Save game history to board (not automatic)?
  * Save as two coordinates, with piece specifier

# User Interface

## clichess
- A commandline chess program for debugging and watching random games.

## browser-sql
- A server api for playing a game and saving it to a sqlite database.

## growser
- A server api using gorilla/websocket for live network chess playing!
- Dependency: Gorilla-websocket (BSD) and Chessboard.js (MIT)
- TODO: Make move with ParseStand() and use chessboardjs's source and newLocation var
  - See http://chessboardjs.com/examples#4003
  - js onDragMove gets tracked...

### Bugs

- Surely, something?

### License

Copyright (C) 2016 Fenimore Love

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

