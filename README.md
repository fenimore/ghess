[![GitHub issues](https://img.shields.io/github/issues/polypmer/go-chess.svg)](https://github.com/polypmer/go-chess/issues) [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/polypmer/go-chess) [![Go Report Card](https://goreportcard.com/badge/github.com/polypmer/go-chess)](https://goreportcard.com/report/github.com/polypmer/go-chess)

# go-chess || ghess

A Golang chess engine and user interfaces. Challenge the **AI**, [Ghess](https://ghess-chess.herokuapp.com/), online.

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
`ghess.go` is a Go package responsible for representing a chess board, parsing PGN input, and validating moves. In the `ui/` directory, `clichess.go` is a simple interface for debugging and `growser.go` is a browser client powered by *websockets*. The `ghess.go` package is broken into `parse.go`, `validation.go`, and `evaluation.go`. See godoc for the [docs](https://godoc.org/github.com/polypmer/go-chess).

- After putting the source in `$GOPATH/src/github.com/polypmer/ghess/`, try

    `go run ui/clichess.go`

- To see a `math/rand` vs `math/rand` game, enter into the **clichess** client:

    `> /random-game`

- To play a real-time game over the internal network, run `growser.go` within the `/ui` directory, and connect to port 8080.

- To use the package in a project, start a new game by calling `ghess.NewBoard()`, which returns an instance of the `ghess.Board` `struct`, ready to `Board.ParseMove()` and return *FEN* positions, `Board.Position()`.

- To evaluate a board position, with positive numbers as a White advantage and negative as Black advantage:

    `> /eval`


# Basic Features and Functionality
- *Most* rules are implemented:
  * Pawns *only* promote to Queen.
  * There is no PGN disambiguation for Queens.
- PGN import-export via `Board.LoadPgn()` and `Board.PgnString()`
- FEN import-export via `Board.LoadFen()` and `Board.Position()`
- Command Line interface.
- Web interface
- Artificial intelligence, goes 5 ply in a few seconds.

# Search and Evaluate Features

- Looks for all possible and valid moves via `Board.SearchForValid()`, which returns two `[]int` slices with the coordinates of possible origins and possible targets. The `Board` field `pieceMap` is a `map[int]string`; the aforementioned `int`s are keys for the standard notation coordinates.
- Evaluation returns a score with a positive value for white advantage and negative value for black advantage. See the `evaluation.go` file for it's emerging api. There is also a `Board.MoveRandom()` method which passes in two `[]int` slices and `math/rand` chooses a move.

----

## Bitmap

The chess engine works with a 120 (10x12) bitmap `[]byte` slice, stored in the `Board` `board` field. This boils down to (accessible with the `/coordinates` command in `clichess.go`):

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

- 11 - 18 1st rank
- 81 - 88 8th rank

----

## Road Map:

1. FIXME: SearchForValid empassant
4. FIXME: Horizon Effect
 - Keep track of capture state to combat horizon effect
4. FIXME: Invalid fen when first Move number is not zero
6. TODO: Change `Board` to `Game`
7. TODO: Save history


### Basic Functionality

- [ ] Minor pawn promotion.
- [ ] Queen disambiguation.
- [ ] Checkmate should update PGN headers/history.
- [ ] `ParseMove` should allow for resign.

----

# User Interfaces

## Clichess
- A commandline chess program, it can output and parse PGN and FEN notation.
- Type `> /help` to list options.

![clichess](http://polypmer.github.io/img/clichess.png "clichess screenshot")


## Growser
- A server api using `gorilla/websocket` for live network chess playing!
- Dependency: gorilla/websocket (BSD) and Chessboard.js (MIT)
- NB. Castling is only when King steps on Rook, not like normals.
- Games are stored with a BoltDB keystore database
- See the repository: [Shallow-Green](https://github.com/polypmer/shallow-green)


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
