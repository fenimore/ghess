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
- Web interface

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

2. Tweak evaluation a bit
4. Deal with Horizon Effect
4. Invalid fen when first number is not zero
5. Add Difficulties to AI UI
4. Keep track of capture state to combat horizon effect
4. Undo the ToUpper and ToLower cause that's not good
6. Make piece map global
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
- A commandline chess program, it can output and parse PGN and FEN notation.
- Type `> /help` to list options.

![clichess](http://polypmer.github.io/img/clichess.png "clichess screenshot")


## Growser
- A server api using `gorilla/websocket` for live network chess playing!
- Dependency: gorilla/websocket (BSD) and Chessboard.js (MIT)
- NB. Castling is only when King steps on Rook, not like normals.
- Games are stored with a BoltDB keystore database
- See the repository: [Shallow-Green](https://github.com/polypmer/shallow-green)

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


After I change the []byte slice board to a [120]byte array, and don't copy it:

    BenchmarkSearchValid-4                       100      22877600 ns/op
    BenchmarkSearchValidSlow-4                    50      29033893 ns/op
    BenchmarkMidGamePruningDepth2-4               10     136505438 ns/op
    BenchmarkOpeningPruningDepth2-4               20      77758483 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1257017288 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2254520731 ns/op
    BenchmarkMidGameTwoPruningDepth3-4        300000          6268 ns/op
    BenchmarkOpeningOrderedDepth3-4                1	1341534583 ns/op
    BenchmarkMidGameOrderedDepth3-4                1	2325314282 ns/op
    BenchmarkMidGameTwoOrderedDepth3-4        200000          6107 ns/op
    BenchmarkOpeningPruningDepth4-4                1	15881901832 ns/op
    BenchmarkMidGamePruningDepth4-4                1	18561026485 ns/op
    PASS

Benchmarks after I figured out that I wasn't calling MiniMaxOrdered inside of MinimaxOrdered...

    BenchmarkMidGamePruningDepth2-4               10     136375546 ns/op
    BenchmarkOpeningPruningDepth2-4               20      91753078 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1322371490 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2164638763 ns/op
    BenchmarkMidGameTwoPruningDepth3-4        300000          6137 ns/op
    BenchmarkOpeningOrderedDepth3-4                1	1246987176 ns/op
    BenchmarkMidGameOrderedDepth3-4                1	2455577971 ns/op
    BenchmarkMidGameTwoOrderedDepth3-4        300000          6135 ns/op
    BenchmarkOpeningPruningDepth4-4                1	15661720638 ns/op
    BenchmarkMidGamePruningDepth4-4                1	18284754487 ns/op

Giving up with PV:

    BenchmarkMidGamePruningDepth2-4               10     131262359 ns/op
    BenchmarkOpeningPruningDepth2-4               20      94373512 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1350511232 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2508213115 ns/op
    BenchmarkOpeningPruningDepth4-4                1	16827821614 ns/op
    BenchmarkMidGamePruningDepth4-4                1	15570438668 ns/op
    PASS
    ok      github.com/polypmer/ghess	50.129s

After Using profiling

    BenchmarkSearchValid-4                       100      15393451 ns/op
    BenchmarkSearchValidSlow-4                   100      15275635 ns/op
    BenchmarkMidGamePruningDepth2-4               20      93899457 ns/op
    BenchmarkOpeningPruningDepth2-4               20      56920214 ns/op
    BenchmarkOpeningPruningDepth3-4                2     784527542 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1518414649 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             3     361195704 ns/op
    BenchmarkOpeningPruningDepth4-4                1	10112804027 ns/op
    BenchmarkMidGamePruningDepth4-4                1	9730502681 ns/op
    PASS
    ok      github.com/polypmer/ghess	31.456s


Benchmarks with pawn validation in Check reduced...

    BenchmarkSearchValid-4                       100      15480015 ns/op
    BenchmarkSearchValidSlow-4                   100      15150476 ns/op
    BenchmarkMidGamePruningDepth2-4               20      91290571 ns/op
    BenchmarkOpeningPruningDepth2-4               30      50329151 ns/op
    BenchmarkOpeningPruningDepth3-4                2     686627412 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1249006486 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             5     342692638 ns/op
    BenchmarkOpeningPruningDepth4-4                1	8347604470 ns/op
    BenchmarkMidGamePruningDepth4-4                1	8536214158 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.632s

Benchmark After certain profiling:

    BenchmarkSearchValid-4                       100      12357042 ns/op
    BenchmarkSearchValidSlow-4                   100      14710691 ns/op
    BenchmarkMidGamePruningDepth2-4               20      78253547 ns/op
    BenchmarkOpeningPruningDepth2-4               30      44239677 ns/op
    BenchmarkOpeningPruningDepth3-4                2     617343700 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1478115205 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             3     342789913 ns/op
    BenchmarkOpeningPruningDepth4-4                1	9527661212 ns/op
    BenchmarkMidGamePruningDepth4-4                1	10551483526 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.863s

Without Updating Check within the Move method

    BenchmarkSearchValid-4                       100      10165640 ns/op
    BenchmarkSearchValidSlow-4                   100      10950169 ns/op
    BenchmarkMidGamePruningDepth2-4               20      55903070 ns/op
    BenchmarkOpeningPruningDepth2-4               50      36001671 ns/op
    BenchmarkOpeningPruningDepth3-4                2     500474578 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1110722333 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             5     239901689 ns/op
    BenchmarkOpeningPruningDepth4-4                1	6221718962 ns/op
    BenchmarkMidGamePruningDepth4-4                1	6580892546 ns/op
    PASS
    ok      github.com/polypmer/ghess	22.101s

Before New Check Method:

    BenchmarkSearchValid-4                      5000        470955 ns/op
    BenchmarkSearchValidSlow-4                  1000       1380257 ns/op
    BenchmarkMidGamePruningDepth2-4               50      42497511 ns/op
    BenchmarkOpeningPruningDepth2-4              100      26641277 ns/op
    BenchmarkOpeningPruningDepth3-4               10     246942839 ns/op
    BenchmarkMidGamePruningDepth3-4                5     279948987 ns/op
    BenchmarkMidGameTwoPruningDepth3-4            20     133189782 ns/op
    BenchmarkOpeningPruningDepth4-4                1	4016534087 ns/op
    BenchmarkMidGamePruningDepth4-4                1	3511034720 ns/op
    PASS
    ok      github.com/polypmer/ghess	24.323s

After new check method

    BenchmarkSearchValid-4                     10000        290186 ns/op
    BenchmarkSearchValidSlow-4                 10000        305094 ns/op
    BenchmarkMidGamePruningDepth2-4               50      21060368 ns/op
    BenchmarkOpeningPruningDepth2-4              100      16300874 ns/op
    BenchmarkOpeningPruningDepth3-4               10     156043108 ns/op
    BenchmarkMidGamePruningDepth3-4               20     146742071 ns/op
    BenchmarkMidGameTwoPruningDepth3-4            20      60358058 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2920717586 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2233602005 ns/op
    PASS
    ok      github.com/polypmer/ghess	19.966s

After Clean Up Tests:

    BenchmarkMidGamePruningDepth2-4              100      25309438 ns/op
    BenchmarkOpeningPruningDepth2-4              100      17696422 ns/op
    BenchmarkOpeningPruningDepth3-4               10     183471772 ns/op
    BenchmarkMidGamePruningDepth3-4               10     147649582 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     623692681 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2224223811 ns/op
    BenchmarkMidGamePruningDepth4-4                1	1748469646 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	5833172254 ns/op
    BenchmarkOpeningPruningDepth5-4                1	25067368045 ns/op
    BenchmarkMidGamePruningDepth5-4                1	12097247498 ns/op
    PASS
    ok      github.com/polypmer/ghess	56.751s


After Checking for Checkmate in Move Function:

    BenchmarkMidGamePruningDepth2-4              100      22388385 ns/op
    BenchmarkOpeningPruningDepth2-4              100      16043151 ns/op
    BenchmarkOpeningPruningDepth3-4                5     243147499 ns/op
    BenchmarkMidGamePruningDepth3-4                3     480370920 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     940100309 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2995183075 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2833983087 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	11650210041 ns/op
    BenchmarkOpeningPruningDepth5-4                1	45387620618 ns/op
    BenchmarkMidGamePruningDepth5-4                1	41308541742 ns/op
    PASS
    ok      github.com/polypmer/ghess	116.066s


With new checkCheck method:

    BenchmarkMidGamePruningDepth2-4              100      22448028 ns/op
    BenchmarkOpeningPruningDepth2-4              100      15045846 ns/op
    BenchmarkOpeningPruningDepth3-4                5     249524429 ns/op
    BenchmarkMidGamePruningDepth3-4                3     508006507 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     919852959 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2828037344 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2467769132 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	10528120071 ns/op
    BenchmarkOpeningPruningDepth5-4                1	44537287486 ns/op
    BenchmarkMidGamePruningDepth5-4                1	40094846903 ns/op
    PASS
    ok      github.com/polypmer/ghess	112.273s


Removing bytes to Upper in Favor of unicode ToLower (woah big gain)

    BenchmarkMidGamePruningDepth2-4              100      16925310 ns/op
    BenchmarkOpeningPruningDepth2-4              100      12206812 ns/op
    BenchmarkOpeningPruningDepth3-4               10     188508894 ns/op
    BenchmarkMidGamePruningDepth3-4                3     362414598 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     719300465 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2205359134 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2047525391 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8160272334 ns/op
    BenchmarkOpeningPruningDepth5-4                1	35592436104 ns/op
    BenchmarkMidGamePruningDepth5-4                1	36422738754 ns/op
    PASS
    ok      github.com/polypmer/ghess	93.830s

Solving Mate in Three Puzzles

    BenchmarkMidGamePruningDepth2-4              100      17760837 ns/op
    BenchmarkOpeningPruningDepth2-4              100      11871520 ns/op
    BenchmarkOpeningPruningDepth3-4               10     189677602 ns/op
    BenchmarkMidGamePruningDepth3-4                3     423109193 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     763345004 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2255497458 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2156571491 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8961648726 ns/op
    BenchmarkOpeningPruningDepth5-4                1	37199162918 ns/op
    BenchmarkMidGamePruningDepth5-4                1	39992108638 ns/op
    PASS
    ok      github.com/polypmer/ghess	100.381s


Fix isUpper method:

    BenchmarkMidGamePruningDepth2-4              100      18516595 ns/op
    BenchmarkOpeningPruningDepth2-4              100      12221270 ns/op
    BenchmarkOpeningPruningDepth3-4               10     192746436 ns/op
    BenchmarkMidGamePruningDepth3-4                3     396922504 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     914956021 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2334247240 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2202312773 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8984800795 ns/op
    BenchmarkOpeningPruningDepth5-4                1	38791150950 ns/op
    BenchmarkMidGamePruningDepth5-4                1	38760424831 ns/op
    PASS
    ok      github.com/polypmer/ghess	101.368s

Redux SearchValid

    BenchmarkMidGamePruningDepth2-4              300       7301357 ns/op
    BenchmarkOpeningPruningDepth2-4              500       4243673 ns/op
    BenchmarkOpeningPruningDepth3-4               30      43071126 ns/op
    BenchmarkMidGamePruningDepth3-4               30      54478924 ns/op
    BenchmarkMidGamePruningDepth3v2-4             10     197960945 ns/op
    BenchmarkOpeningPruningDepth4-4                2     755648819 ns/op
    BenchmarkMidGamePruningDepth4-4                2     665499051 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	2091657012 ns/op
    BenchmarkOpeningPruningDepth5-4                1	8319886414 ns/op
    BenchmarkMidGamePruningDepth5-4                1	5527034454 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.643s

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
