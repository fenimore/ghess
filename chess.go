/*
Go Chess Engine
Fenimore Love 2016
GPLv3

TODO: Search and Evaluation
TODO: Fen PGN reading
TODO: Fen output
*/
package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"regexp"
	"bufio"
	"os"
)

// The chessboard type
type Board struct {
	board []byte // piece position
	// Game Variables
	castle    []byte // castle possibility KQkq or ----
	empassant int    // square vulnerable to empassant
	score     string
	toMove    string // Next move is w or b
	moves     int    // the count of moves
	// Map for display grid
	pgnMap map[string]int    // the pgn format
	pieces map[string]string // the unicode fonts
	// Game Positions
	fen string
	pgn string
}

// __init__ for Board
func NewBoard() Board {
	b := make([]byte, 120)
	fmt.Println("Initializing new Chess game\n")

	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)

	// Printed Board Notations
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'

	// Map of PGN notation
	m := make(map[string]int)
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"] = 11, 12, 13, 14, 15, 16, 17, 18
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"] = 21, 22, 23, 24, 25, 26, 27, 28
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"] = 31, 32, 33, 34, 35, 36, 37, 38
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"] = 41, 42, 43, 44, 45, 46, 47, 48
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"] = 51, 52, 53, 54, 55, 56, 57, 58
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"] = 61, 62, 63, 64, 65, 66, 67, 68
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"] = 71, 72, 73, 74, 75, 76, 77, 78
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"] = 81, 82, 83, 84, 85, 86, 87, 88

	// Map of unicode fonts
	r := make(map[string]string)
	r["p"], r["P"] = "\u2659", "\u265F"
	r["b"], r["B"] = "\u2657", "\u265D"
	r["n"], r["N"] = "\u2658", "\u265E"
	r["r"], r["R"] = "\u2656", "\u265C"
	r["q"], r["Q"] = "\u2655", "\u265B"
	r["k"], r["K"] = "\u2654", "\u265A"
	r["."] = "\u2022"

	return Board{
		board:  b,
		castle: []byte(`KQkq`),
		pgnMap: m,
		pieces: r,
		toMove: "w",
		score:  "*",
		moves:  1,
	}
}

// Return a string of the board
// Todo Unicode chess pieces
func (b *Board) String() string {
	var printBoard string
	for idx, val := range b.board {
		if idx < 100 && idx > 10 {
			if idx%10 != 0 && idx < 90 {
				if (idx+1)%10 != 0 { // why not || ?
					font := b.pieces[string(val)]
					printBoard += "|" + font + "| "
				} else {
					printBoard += ":" + string(val)
				}
			}
		}
		if idx > 90 && idx < 99 {
			printBoard += ":" + string(val) + ": "
		} else if idx%10 == 0 && idx != 0 {
			printBoard += "\n"
		}
	}
	return printBoard
}

/*
Move and validation
*/
// Wrapper in standard notation
func (b *Board) pgnMove(orig, dest string) error {
	e := b.Move(b.pgnMap[orig], b.pgnMap[dest])
	if e != nil {
		return e
	}
	return nil
}

// Move byte value to new position
func (b *Board) Move(orig, dest int) error {
	fmt.Print("Moves: ", b.moves, " | Castle: ", string(b.castle))
	fmt.Println(" | Turn: ", b.toMove)
	val := b.board[orig]
	var o byte // supposed starting square
	var d byte // supposed destination
	if b.toMove == "w" {
		// check that orig is Upper
		// and dest is Enemy or Empty
		o = []byte(bytes.ToUpper(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToLower(b.board[dest : dest+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		// and dest is Enemy or Empty
		o = []byte(bytes.ToLower(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToUpper(b.board[dest : dest+1]))[0]
	}
	// Check if it is the right turn
	if b.board[orig] != o {
		return errors.New("Not your turn")
	}
	// Check if Origin is Empty
	if o == '.' {
		return errors.New("Empty square")
	}
	// Check if destination is Enemy
	if b.board[dest] != d { //
		return errors.New("Can't attack your own piece")
	}
	p := string(bytes.ToUpper(b.board[orig : orig+1]))
	fmt.Print("Test: ", orig, dest)
	switch {
	case p == "P":
		e := b.validPawn(orig, dest, d)
		if e != nil {
			return e
		}
	case p == "N":
		fmt.Print("is knight")
	case p == "B":
		fmt.Print("is bishop")
	case p == "R":
		fmt.Print("is rook")
	case p == "Q":
		fmt.Print("is queen")
	case p == "K":
		fmt.Print("is king")
	}
	// Update Board
	b.board[orig] = '.'
	b.board[dest] = val // check if it was the King
	// Update Game variables
	if b.toMove == "w" {
		b.toMove = "b"
	} else {
		b.moves++ // add one to move
		b.toMove = "w"
	}
	return nil
}

// validate Pawn Move
func (b *Board) validPawn(orig int, dest int, d byte) error {
	err := errors.New("Illegal Move")
	var remainder int
	if b.toMove == "w" {
		remainder = dest - orig
	} else if b.toMove == "b" {
		remainder = orig - dest
	}
	if remainder == 10 {
		// regular move
	} else if remainder == 20 { // two spaces
		// double starter move
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		}
		if orig < 70 && b.toMove == "b" {
			return err
		}
	} else if remainder == 9 || remainder == 11 {
		// Attack vector
		// check if b.board[orig+10] == '.'
		if b.board[dest] == d && d != '.' {
			// Proper attack
		} else {
			return err
		}
	}
	return nil
}

// Valid Knight
// remainder == 21 19 12 8 -21 -19 -12 -8
// Valid Bishop
// Valid Rook
// Valid Queen
// Valid King
// remainder == 10 11 1 9 -10 -9 -11 -1

/*
TODO: Export fen
TODO: Parse fen
TODO: Parse pgn
*/
// Par example:
// e4 ...

func (b *Board) pgnTest(move string) error {
	re, _ := regexp.Compile(`(.)x(..)`) // want to know what is in front of 'x'
	res := re.FindStringSubmatch(move)
	fmt.Println(res[2])
	return nil
}

func (b *Board) parsePgn(move string) error {
	move = strings.TrimRight(move, "\r\n") // prepare for input
	/*
          Regex Pattern: [B-R]?[a-h]?x?[a-h]\d{1}
          Examples Inputs: 
                   e4 | d5 | exd5 | Bc7 | Qxc7
        */
	var orig int // find origin coord of move
	var square string // find pgnMap key of move
	var attacker string // left of x
	var piece string // find move piece
	var precise string // for multiple possibilities

	// Check if Capture
	isCapture, _ := regexp.MatchString(`x`, move)
	if isCapture && len(move) > 3 {
		re, _ := regexp.Compile(`(.)x(..)`) // want to know what is in front of 'x'
		res := re.FindStringSubmatch(move)
		attacker = res[1]
		if attacker == strings.ToLower(attacker){
			piece = "P"
		}
		square = res[2]
	} else { // No x
		chars := len(move)
		switch {
		case: chars == 2
			// do somethin
			piece = "P"
			square = move
		case: chars == 3:
			piece = move[0]
		case: chars ==  4:
			piece = move[0]
			precise := move[1]
		case: chars < 2:
			return errors.New("Not enough input")
		}

	}
	// 
	dest := b.pgnMap[square]
	piece = strings.ToUpper(piece)
	switch {
	case piece == "P":
		var target1, target2 int
		var t byte
		if b.toMove == "w" {
			t = 'P' // target origin
			if isCapture {
				target1, target2 = dest-9, dest-11
			} else {
				target1, target2 = dest-10, dest-20
			}
		} else {// is black to move
			t = 'p' // target origin
			if isCapture {
				target1, target2 = dest+9, dest+11
			} else {
				target1, target2 = dest+10, dest+20				}
		}
		if b.board[target1] == t {
			orig = target1
		} else if b.board[target2] == t {
			orig = target2
		}
	case piece == "N":
		// do something
	}
	if orig != 0 && dest != 0 {
		err := b.Move(orig, dest)
		return err
	} else {
		return  errors.New("No such move")
	}
}

func (b *Board) parseFen() {
	// Parse Fen
}

func (b *Board) genFen() string {
	// b.board -> Fen
	fen := "Fen string"
	return fen
}

/*
Main thread
*/
func main() {
	board := NewBoard()
	fmt.Println("coordinates:")
	board.Coordinates()
	PlayGame(board)
	//TestGame(board)
	//board.pgnTest("exd5")

}

/*
Helper Testing method
*/

func TestGame(board Board) {
	e := board.parsePgn("e4")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
	e = board.parsePgn("d5")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
	e = board.parsePgn("exd5")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
}

func PlayGame(board Board) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(board.String())
	for {
		fmt.Print(board.toMove, " to move: ")
		input, _ := reader.ReadString('\n')
		e := board.parsePgn(input)
		if e != nil {
			fmt.Print(e)
		}
		fmt.Print(board.String())
	}
}

func (b *Board) Coordinates() {
	for idx, val := range b.board {
		if idx < 100 && idx > 10 {
			if idx%10 != 0 && idx < 90 {
				if (idx+1)%10 != 0 {
					fmt.Print(":", idx, ": ")
				} else {
					fmt.Print(":", string(val))
				}
			}
		}
		if idx > 90 && idx < 99 {
			fmt.Print(": ", string(val), ": ")
		} else if idx%10 == 0 && idx != 0 {
			fmt.Print("\n")
		}
	}
}
