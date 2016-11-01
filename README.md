[![GitHub issues](https://img.shields.io/github/issues/polypmer/go-chess.svg)](https://github.com/polypmer/go-chess/issues) [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/polypmer/go-chess) [![Go Report Card](https://goreportcard.com/badge/github.com/polypmer/go-chess)](https://goreportcard.com/report/github.com/polypmer/go-chess)

# go-chess || ghess

A Golang chess engine and user interfaces

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
- Simple AI

# Search and Evaluate Features

- Looks for all possible and valid moves via `Board.SearchForValid()`, which returns two `[]int` slices with the coordinates of possible origins and possible targets. The `Board` field `pieceMap` is a `map[int]string`; the aforementioned `int`s are keys for the standard notation coordinates.
- Evaluation returns a score with a positive value for white advantage and negative value for black advantage. See the `evaluation.go` file for it's emerging api. There is also a `Board.MoveRandom()` method which passes in two `[]int` slices and `math/rand` chooses a move.
- A MiniMax algorithm is a work in progress, but with it one can play a game against a rather low level opponent.

----

## `Board` struct fields:
- [Game variables]
  * board []byte [the piece positions]
  * castle []byte [Remaining possibilities]
  * score string
  * toMove string [either "w" or "b"]
  * empassant int [coordinate value]
  * checkmate bool
  * check bool
- [Board display]
  * pgnMap map[string]int
  * pieceMap map[int]string
  * pieces map[string]string [unicode fonts]
  * rows map[int][8]int [for white/black coloring]
- moves [move tally]
- fen string[Position string]
- pgn string [move history]
- headers string [for PGN]
- pgnPattern *regexp.Regexp
- fenPattern *regexp.Regexp
- history []int

As long as it remains on my `TODO` list to change, these fields are unexported and for accessing this data, one can call `Board.Stats()` to return a `map[string]string` of various imminently useful fields, such as "turn", "check", and "scores" (see `ghess.go` `Stats()` for a complete list.

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

## TODO General

1. More tests.
4. Export `Board` fields.
6. Change `Board` to `Game`, as that makes more sense...
3. Convert `ParseStand()` to PGN move.
7. Save history

### TODO Basic rules

- Minor pawn promotion.
- Queen disambiguation.

### TODO Basic Functionality

- Checkmate should update PGN headers/history.
- `ParseMove` should allow for resign.

### TODO Extra features

- Move history/ Undo
  * Save game history to board (not automatic)?
  * Save as two coordinates, with piece specifier

----

# User Interfaces

## Clichess
- A commandline chess program for debugging and watching random games.
- Type `> /help` to list options.

![clichess](http://polypmer.github.io/img/clichess.png "clichess screenshot")


## Growser
- A server api using `gorilla/websocket` for live network chess playing!
- Dependency: gorilla/websocket (BSD) and Chessboard.js (MIT)
  * TODO: Add watch random
  * TODO: Add play AI
  * TODO: context and game index...
  * TODO: everything user
- Castling is only when King steps on Rook, not like normals.

## browser-sql
- A server api for playing a game and saving it to a sqlite database.


----

# Benchmarks:

(These are for my testing purposes). As of Nov 1, using the suped up `SearchValid()`, Minimax Benchmarks look like this:


    BenchmarkSearchValid-4                    50      25884559 ns/op
    BenchmarkSearchValidSlow-4                50      27821677 ns/op
    BenchmarkMidGamePruningDepth2-4           10     138513378 ns/op
    BenchmarkOpeningPruningDepth2-4           20      94688395 ns/op
    BenchmarkOpeningPruningDepth3-4            1	1441998152 ns/op
    BenchmarkMidGamePruningDepth3-4            1	2177541792 ns/op
    BenchmarkOpeningPruningDepth4-4            1	16566151366 ns/op
    BenchmarkMidGamePruningDepth4-4            1	16079072907 ns/op



### Bugs

- See issues.

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
