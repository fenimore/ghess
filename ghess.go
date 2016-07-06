/*
Go Chess Engine - Ghess
Fenimore Love 2016
GPLv3

TODO: Search and Evaluation
TODO: Fen PGN reading
TODO: Fen output
*/
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"strconv"
	"time"
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
	check     bool
	// Map for display grid
	pgnMap map[string]int    // the pgn format
	pieceMap map[int] string // coord to standard notation
	pieces map[string]string // the unicode fonts
	// Game Positions
	fen      string
	pgn      string
	pgnHeaders  string
	pattern *regexp.Regexp // For parsing PGN

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
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"] = 18, 17, 16, 15, 14, 13, 12, 11
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"] = 28, 27, 26, 25, 24, 23, 22, 21
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"] = 38, 37, 36, 35, 34, 33, 32, 31
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"] = 48, 47, 46, 45, 44, 43, 42, 41
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"] = 58, 57, 56, 55, 54, 53, 52, 51
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"] = 68, 67, 66, 65, 64, 63, 62, 61
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"] = 78, 77, 76, 75, 74, 73, 72, 71
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"] = 88, 87, 86, 85, 84, 83, 82, 81

	// Todo make map for pieceMap[]
	
	// Map of unicode fonts
	r := make(map[string]string)
	r["p"], r["P"] = "\u2659", "\u265F"
	r["b"], r["B"] = "\u2657", "\u265D"
	r["n"], r["N"] = "\u2658", "\u265E"
	r["r"], r["R"] = "\u2656", "\u265C"
	r["q"], r["Q"] = "\u2655", "\u265B"
	r["k"], r["K"] = "\u2654", "\u265A"
	r["."] = "\u00B7"

	pattern,_ := regexp.Compile(`([PNBRQK]?[a-h]?[1-8]?)x?([a-h][1-8])([\+\?\!]?)|O(-?O){1,2}`)
	return Board{
		board:  b,
		castle: []byte(`KQkq`),
		pgnMap: m,
		pieces: r,
		toMove: "w",
		score:  "*",
		moves:  1,
		pattern: pattern,
	}
}

// Set pgnHeaders for a pgn export
func (b *Board) setHeaders(w, bl string) {
	w = strings.TrimRight(w, "\r\n")
	bl = strings.TrimRight(bl, "\r\n")
	y, m, d := time.Now().Date()
	ye, mo, da := strconv.Itoa(y), strconv.Itoa(int(m)),
	strconv.Itoa(d)
	white  := "[White \""+w+"\"]"
	black  := "[Black \""+bl+"\"]"
	date   := "[Date \""+ye+"."+mo+"."+da+"\"]"
	result := `[Result "*"]` 
	b.pgnHeaders = white + "\n" + black+ "\n" + date + "\n" + result + "\n"
}

// Return a string of the board
func (b *Board) RotateBroken() string {
	// This is Broken... Yikes
	var printBoard string
	for idx, val := range b.board {
		if idx < 100 && idx > 10 {
			if idx%10 != 0 && idx < 90 {
				if (idx+1)%10 != 0 { // why not || ?
					font := b.pieces[string(val)]
					printBoard += "|" + font + "|"
				} else {
					printBoard += ":" + string(val)
				}
			}
		}
		if idx > 90 && idx < 99 {
			printBoard += ":" + string(val) + ":"
		} else if idx%10 == 0 && idx != 0 {
			printBoard += "\n"
		}
	}
	return printBoard
}

func (b *Board) String() string {
	// TODO Rotate Board
	game := b.board
	p := b.pieces
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0  {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i])+": "
			continue
		}
		printBoard += "|"+p[string(game[i])]+"|"
	}

	printBoard += "\n"
	printBoard += "   :a::b::c::d::e::f::g::h:\n"
	return printBoard
}

/*
Move and validation
*/
// Wrapper in portable game notation
// 'Two' coordinate notation
func (b *Board) pgnMove(orig, dest string) error {
	//e2e4
	e := b.Move(b.pgnMap[orig], b.pgnMap[dest])
	if e != nil {
		return e
	}
	return nil
}

// Move byte value to new position
func (b *Board) Move(orig, dest int) error {
	val := b.board[orig]
	var o byte // supposed starting square
	var d byte // supposed destination
	var empassant bool //refactor?
	var isCastle bool

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
	// Check for Castle
	if orig == 14 {
		isCastle = b.board[dest] == 'R'
	} else if orig == 84 {
		isCastle = b.board[dest] == 'r'
	}

	err := b.basicValidation(orig, dest, o, d, isCastle)
	if err != nil {
		return err
	}
	
	p := string(bytes.ToUpper(b.board[orig : orig+1]))
	switch {
	case p == "P":
		e := b.validPawn(orig, dest, d)
		if e != nil {
			return e
		}
		emp := dest - orig
		if emp > 11 || emp < -11 {
			empassant = true
		}
	case p == "N":
		e := b.validKnight(orig, dest)
		if e != nil {
			return e
		}
	case p == "B":
		e := b.validBishop(orig, dest)
		if e != nil {
			return e
		}
	case p == "R":
		e := b.validRook(orig, dest)
		if e != nil { 
			return e
		}
		switch { // Castle
		case orig == b.pgnMap["a1"]:
			b.castle[1] = '-'
		case orig == b.pgnMap["a8"]:
			b.castle[3] = '-'
		case orig == b.pgnMap["h1"]:
			b.castle[0] = '-'
		case orig == b.pgnMap["h8"]:
			b.castle[2] = '-'
		}
	case p == "Q":
		e := b.validQueen(orig, dest)
		if e != nil {
			return e
		}
	case p == "K": // is castle?
		if !isCastle {
			e := b.validKing(orig, dest, false)
			if e != nil {
				return e
			}
			if orig == 14 || orig == 84 { // starting pos
				switch {
				case o == 'K':
					b.castle[0], b.castle[1] = '-', '-'
				case o == 'k':
					b.castle[2], b.castle[3] = '-', '-'
				}
			}
		} else {
			e := b.validKing(orig, dest, true)
			if e != nil {
				return e
			}
			
		}
	}
	// Update Board
	b.board[orig] = '.'
	if !isCastle {
		b.board[dest] = val
	} else { // castle
		if dest > orig { // queen side
			b.board[dest-2],
			b.board[dest-3] = val, b.board[dest]
		} else {         // king side
			b.board[dest+1],
			b.board[dest+2] = val, b.board[dest]
		}
		b.board[dest] = '.'
	}
	// TODO check for Check
	// Update Game variables
	if b.toMove == "w" {
		b.toMove = "b"
	} else {
		b.moves++ // add one to move count
		b.toMove = "w"
	}
	if empassant {
		b.empassant = dest
	} else {
		b.empassant = 0
	}
	return nil
}

func (b *Board) basicValidation(orig, dest int, o, d byte, isCastle bool) error {
	// Check if it is the right turn
	if b.board[orig] != o {
		return errors.New("Not your turn")
	}
	// Check if Origin is Empty
	if o == '.' {
		return errors.New("Empty square")
	}
	// Check if destination is Enemy
	if b.board[dest] != d && !isCastle { //
		return errors.New("Can't attack your own piece")
	}
	return nil
}

// validate Pawn Move
func (b *Board) validPawn(orig int, dest int, d byte) error {
	err := errors.New("Illegal Pawn Move")
	var remainder int
	var empOffset int
	var empTarget byte
	if b.toMove == "w" {
		remainder = dest - orig
		empOffset = -10 // where the empassant piece should be
		empTarget = 'p'
	} else if b.toMove == "b" {
		remainder = orig - dest
		empOffset = 10
		empTarget = 'P'
	}
	if remainder == 10 {
		// regular move
	} else if remainder == 20 { // two spaces
		// double starter move
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		} else if orig < 70 && b.toMove == "b" {
			return err
		}
	} else if remainder == 9 || remainder == 11 {
		// Attack vector
		// check if b.board[dest+10] == '.'
		if b.board[dest] == d && d != '.' {
			// Proper attack
		} else if b.board[dest] == d && dest+empOffset == b.empassant{
			// Empassant attack
			if b.board[dest+empOffset] == empTarget { // is the right case
				b.board[b.empassant] = '.'
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (b *Board) validKnight(orig int, dest int) error {
	var possibilities [8]int
	possibilities[0], possibilities[1],
	possibilities[2], possibilities[3],
	possibilities[4], possibilities[5],
	possibilities[6], possibilities[7] = orig+21,
	orig+19, orig+12, orig+8, orig-8,
	orig-12, orig-19, orig-21

	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	return errors.New("Illegal Knight Move")
}

func (b *Board) validBishop(orig int, dest int) error {
	// Check if other pieces are in the way
	err := errors.New("Illegal Bishop Move")
	trajectory := orig - dest
	a1h8 := trajectory % 11 // if 0 remainder...
	a8h1 := trajectory % 9
	// Check which slope
	if a1h8 == 0 {
		if dest > orig { // go to bottom right
			for i := orig + 11; i <= dest-11; i += 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if dest < orig { // go to top left
			for i := orig - 11; i >= dest+11; i -= 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if a8h1 == 0 {
		if dest > orig { // go to bottem left
			for i := orig + 9; i <= dest-9; i += 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if orig > dest { // go to top right
			for i := orig - 9; i >= dest+9; i -= 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	}
	return nil
}

func (b *Board) validRook(orig int, dest int) error {
	// Check if pieces are in the way
	err := errors.New("Illegal Rook Move")
	fmt.Println(dest, orig)
	fmt.Println(string(b.board[15]))
	remainder := dest - orig
	if remainder < 10 && remainder > -10 {
		// Horizontal
		if remainder < 0 {
			for i := orig - 1; i > dest; i-- {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 1; i < dest; i++ {
				if b.board[i] != '.' {
					fmt.Println("i is ", i)
					return err
				}
			}
		}
	} else {
		// Vertical
		if remainder < 0 {
			for i := orig - 10; i > dest; i -= 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 10; i < dest; i += 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	}
	return nil
}

func (b *Board) validQueen(orig int, dest int) error {
	remainder := dest - orig
	vertical := remainder%10 == 0
	horizontal := remainder < 9 && remainder > -9 // Horizontal
	diagA8 := remainder%9 == 0                    // Diag a8h1
	diagA1 := remainder%11 == 0                   // Diag a1h8
	if horizontal {                               // should be first?
		fmt.Println("Horizontal")
	} else if vertical { // then it doesn't matter
		fmt.Println("Vertical")
	} else if diagA8 {
		fmt.Println("Diag")
	} else if diagA1 {
		fmt.Println("Diag")
	} else {
		return errors.New("Illegal Queen Move")
	}
	// check if anything is inbetween

	return nil
}

// do castle in King validation
func (b *Board) validKing(orig int, dest int, castle bool) error {
	castlerr := errors.New("Something is in your way")
	noCastle := errors.New("Castle on this side is foutu")
	var possibilities [8]int
	g := b.board // g for gameboard
	possibilities[0], possibilities[1],
	possibilities[2], possibilities[3],
	possibilities[4], possibilities[5],
	possibilities[6], possibilities[7] = orig+10,
	orig+11, orig+1, orig+9, orig-10,
	orig-11, orig-1, orig-9
	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	if castle {
		queenSideCastle := !(g[orig+1] != '.' || g[orig+2] != '.' || g[orig+3] != '.')
		kingSideCastle := !(g[orig-1] != '.' || g[orig-2] != '.')
		if dest > orig { // Queen side
			if !queenSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[1] != 'Q'{
					return noCastle
				}
				b.castle[0], b.castle[1] = '-','-'
			} else { // b
				if b.castle[3] != 'q' {
					return noCastle 
				}
				b.castle[2], b.castle[3] = '-','-'
			}
		} else if orig > dest { 
			if !kingSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[0] != 'K'{
					return noCastle
				}
				b.castle[0], b.castle[1] = '-','-'
			} else {
				if b.castle[2] != 'k' {
					return noCastle 
				}
				b.castle[2], b.castle[3] = '-','-'
			}
		}

	} else {
		return errors.New("Illegal King Move")
	}
	return nil
}



/*
TODO: Export fen
TODO: Parse fen
TODO: Parse pgn
Pgn parse:
  Accept check/checkmate indicaters
  Implement specific pieces..
  Dont all taking a piece from simple moving
*/

func (b *Board) ParsePgn(move string) error {
	move = strings.TrimRight(move, "\r\n") // prepare for input
	// Variables
	var piece string    // find move piece	
	var orig int        // find origin coord of move
	var square string   // find pgnMap key of move
	var attacker string // left of x
	//var precise string // for multiple possibilities
	var target byte // the piece to move, in proper case

	// Status
	isCastle := false
	isWhite := b.toMove == "w"
	isCapture, _ := regexp.MatchString(`x`, move)
	
	res := b.pattern.FindStringSubmatch(move)
	if res == nil && move != "O-O" && move != "O-O-O" {
		return errors.New("invalid input")
	} else if move == "O-O" || move == "O-O-O" {
		// be nice
		isCastle = true
	} 
	
	// Either is catpure or not
	if isCapture {
		attacker = res[1]
		if attacker == strings.ToLower(attacker) {
			piece = "P"
		} else { // if  upper case, forcement a piece
			piece = res[1]
		}
		square = res[2]
	} else if isCastle {
		if move == "O-O" {
			piece = "K"
			if isWhite {
				square = "h1"
			} else {
				square = "h8"
			}
		} else if move == "O-O-O" {
			piece = "K"
			if isWhite {
				square = "a1"
			} else {
				square = "a8"
			}
		}
	} else { // No x
		chars := len(move)
		if chars == 2 {
			piece = "P"
			square = res[2]
		} else if chars == 3 && move != "0-0" {
			// Breaks when e44 is entered...
			piece = res[1]
			square = res[2] //move[0]
		} else if chars == 4 {
			piece = res[1] // remove second char
			//precise = move
			square = res[2]
		} else {
			return errors.New("Not enough input")
		}
	}
	// the presumed destination
	dest := b.pgnMap[square]
	// The piece will be saved as case sensitive byte
	if b.toMove == "b" {
		target = []byte(strings.ToLower(piece))[0]
	} else {
		target = []byte(piece)[0]
	}
	switch {
	case piece == "P": // Pawn Parse
		var possibilities [2]int // two potentional origins
		// TODO: Allow for empassant take
		if b.toMove == "w" {
			if isCapture {
				possibilities[0],
					possibilities[1] = dest-9,
					dest-11
			} else {
				possibilities[0],
					possibilities[1] = dest-10,
					dest-20
			}
		} else { // is black to move
			if isCapture {
				possibilities[0],
					possibilities[1] = dest+9,
					dest+11
			} else {
				possibilities[0],
					possibilities[1] = dest+10,
					dest+20
			}
		}
		if b.board[possibilities[0]] == target {
			orig = possibilities[0]
		} else if b.board[possibilities[1]] == target {
			orig = possibilities[1]
		}
	case piece == "N": // Knight Parse
		var possibilities [8]int
		// TODO: assume no precision
		// Change to possibilities[]
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+21,
			dest+19, dest+12, dest+8, dest-8,
			dest-12, dest-19, dest-21
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	case piece == "B": // Bishop Parse
		var possibilities [14]int
		ticker := 0
		// a8 - h1
		for i := dest + 9; i < 90; i += 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 9; i > 10; i -= 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// a1 - h8 Vector
		for i := dest + 11; i < 90; i += 11 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 11; i > 10; i -= 11 {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// Find piece origin
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	case piece == "R": // Rook Parse
		var possibilities [14]int
		ticker := 0
		// Horizontal Vector
		for i := dest + 10; i < 90; i += 10 {
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 10; i > 10; i -= 10 {
			possibilities[ticker] = i
			ticker++
		}
		// Vertical Vector
		for i := dest + 1; i < 90; i++ {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 1; i > 10; i-- {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
	Looposs:
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				fmt.Println(possibility, string(target))
				orig = possibility
				err := b.validRook(orig, dest)
				if err != nil {
					continue
				}
				break Looposs
			}
		}
	case piece == "Q": // Queen Parse
		for idx, possibility := range b.board {
			if possibility == target {
				orig = idx
				break
			}
		}
	case piece == "K": // King Parse
		var possibilities [8]int
		if isCastle {
			if isWhite {
				orig = 14
			} else {
				orig = 84
			}
			break
		} 
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+10,
			dest+11, dest+1, dest+9, dest-10,
			dest-11, dest-1, dest-9
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	}
	// Move the Piece
	// - Validate Move in Board.Move()
	if b.board[dest] != '.' && !isCapture && !isCastle {
		fmt.Print(isCastle)
		return errors.New("Not the proper capture syntax")

	}
	if orig != 0 && dest != 0 {
		err := b.Move(orig, dest)
		if err == nil {
			// Update pgn History
			if b.toMove == "b"{
				b.pgn += strconv.Itoa(b.moves) +". "
			}
			b.pgn += (move + " ")
		}
		return err
	} else {
		return errors.New("No such move")
	}
}

// Read a Pgn match
func (b *Board) readPgnMatch(match string) (Board, error) {	
	game := NewBoard()
	result := game.pattern.FindAllString(match, -1)
	for _, val := range result {
		fmt.Print(game.String(), game.moves)
		err := game.ParsePgn(val)
		if err != nil {
			return game, err
		}
	}
	// if error, could not read
	return game, nil
}

func (b *Board) stringPgn() string {
	
	return b.pgn
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
	PlayGame(board)
	//TestGame(board)
}

/*
Helper Testing method
*/

func TestGame(board Board) {
	e := board.ParsePgn("e4")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
}

func PlayGame(board Board) { // TODO Rotate Board
	var turn string
	welcome := `
/~ |_ _  _ _
\_|||(/__\_\




go-chess
    Enter /help for more options

`
	manuel := `Help:
    Prefix commands with / - slash

Commands:
	quit - exit game
	new - new game
	coordinates - print board coordinates
	pgn - print PGN history
	fen - print FEN position
	set-headers - set PGN headers
	headers - print game info
Tests:
	test-castle - test castling`
	reader := bufio.NewReader(os.Stdin)
	// welcome message
	fmt.Println(welcome)
	fmt.Print(board.String())
	
	for {
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Print(turn, " to move: ")
		input, _ := reader.ReadString('\n')
		isCmd, _ := regexp.MatchString(`/`, input)
		if isCmd {
			input = strings.TrimRight(input, "\r\n")
			switch {
			case input == "/help":
				fmt.Print("\n", manuel)
			case input == "/quit":
				os.Exit(1)
			case input == "/new":
				board = NewBoard()
				fmt.Print(board.String())
			case input == "/print":
				fmt.Print(board.String())
			case input == "/coordinates":
				fmt.Println("Coordinates:")
				board.Coordinates()
			case input == "/pgn":
				fmt.Println("PGN history:")
				fmt.Println(board.pgn, "\n")
			case input == "/read-pgn":
				var err error
				fmt.Print("Enter PGN history: ")
				history, _ := reader.ReadString('\n')
				board, err = board.readPgnMatch(history)
				if err != nil {
					fmt.Println(err)
				}
			case input == "/fen":
				fmt.Println("FEN position:")
			case input == "/set-headers":
				fmt.Println("Set Headers:")
				fmt.Print("White: ")
				inWhite, _ := reader.ReadString('\n')
				fmt.Print("Black: ")
				inBlack, _:= reader.ReadString('\n')
				board.setHeaders(inWhite, inBlack)
			case input == "/headers":
				fmt.Println(board.pgnHeaders)
			case input == "/test-castle":
				board = NewBoard()
				board = TestCastle(board)
				fmt.Print(board.String())
			case input == "/test-pawn":
				board = NewBoard()
				board = TestPawn(board)
				fmt.Print(board.String())
			case input == "/test-pgn":
				hist := `1. b4 g6 2. c4 Nf6 3. Bb2 Bg7 4. Qc2 Nc6 5. Nc3 b6 6. Nf3 Bb7 7. d4 d5 8. g3 Qd7`
				var err error
				board, err = board.readPgnMatch(hist)
				if err != nil {
					fmt.Println(err)
					board = NewBoard()
				}
				fmt.Print(board.String())
			default:
				fmt.Println("Mysterious input")
			}

			continue
		}
		e := board.ParsePgn(input)
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Println("\n-------------------")
		fmt.Print("Debug:\nMove: ", board.moves,
			" | Castle: ", string(board.castle))
		fmt.Println(" | Turn: ", turn)
		if e != nil {
			fmt.Printf("   [Error: %v]\n", e)
		}
		fmt.Print(board.String())
	}
}

func TestCastle(board Board) Board {
	// Castle Test
	_ = board.ParsePgn("b4")
	_ = board.ParsePgn("g6")
	_ = board.ParsePgn("c4")
	_ = board.ParsePgn("Nf6")
	_ = board.ParsePgn("Bb2")
	_ = board.ParsePgn("Bg7")
	_ = board.ParsePgn("Qc2")
	_ = board.ParsePgn("Nc6")
	_ = board.ParsePgn("Nc3")
	_ = board.ParsePgn("b6")
	_ = board.ParsePgn("Nf3")
	_ = board.ParsePgn("Bb7")
	_ = board.ParsePgn("d4")
	_ = board.ParsePgn("d5")	
	_ = board.ParsePgn("g3")
	_ = board.ParsePgn("Qd7")
	return board
}

func TestPawn(board Board) Board {
	// Castle Test
	_ = board.ParsePgn("e4")
	_ = board.ParsePgn("d5")
	_ = board.ParsePgn("exd5")
	return board
}

// Print the Board.Move() coordinate
func (b *Board) Coordinates() {
		// TODO Rotate Board
	game := b.board
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0  {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i])+": "
			continue
		}
		printBoard += "|"+strconv.Itoa(i)+"|"
	}

	printBoard += "\n"
	printBoard += "   :a ::b ::c ::d ::e ::f ::g ::h :\n"
	fmt.Println(printBoard)
}
