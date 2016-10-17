// Package ghess is a Go Chess Engine.
//
// Fenimore Love 2016
// GPLv3
//
package ghess

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Board is a chessboard type.
// TODO: Make Upper Case? M-c for upper case
type Board struct {
	board []byte // piece position
	// Game Variables
	castle    []byte // castle possibility KQkq or ----
	empassant int    // square vulnerable to empassant
	score     string
	toMove    string // Next move is w or b
	moves     int    // the count of moves
	Check     bool
	Checkmate bool // start Capitalizing
	// Map for display grid
	pgnMap   map[string]int    // the pgn format
	pieceMap map[int]string    // coord to standard notation
	PieceMap map[int]string    // exported cord to stand notation
	pieces   map[string]string // the unicode fonts
	rows     map[int][8]int    // rows for white/black squaring
	// Game Positions
	fen        string         // Game position
	pgn        string         // Game history
	headers    string         // Pgn format
	pgnPattern *regexp.Regexp // For parsing PGN
	fenPattern *regexp.Regexp // For validating FEN input
	history    [8]int         // For Draws, last six coordinates
}

// NewBoard returns pointer to new Board in the starting position.
func NewBoard() Board {
	b := make([]byte, 120)
	dictionary()
	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)

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
	// pieceMap
	p := make(map[int]string)
	p[18], p[17], p[16], p[15], p[14], p[13], p[12], p[11] = "a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1"
	p[28], p[27], p[26], p[25], p[24], p[23], p[22], p[21] = "a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2"
	p[38], p[37], p[36], p[35], p[34], p[33], p[32], p[31] = "a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3"
	p[48], p[47], p[46], p[45], p[44], p[43], p[42], p[41] = "a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4"
	p[58], p[57], p[56], p[55], p[54], p[53], p[52], p[51] = "a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5"
	p[68], p[67], p[66], p[65], p[64], p[63], p[62], p[61] = "a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6"
	p[78], p[77], p[76], p[75], p[74], p[73], p[72], p[71] = "a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7"
	p[88], p[87], p[86], p[85], p[84], p[83], p[82], p[81] = "a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8"
	// Map of unicode fonts
	r := make(map[string]string)
	r["p"], r["P"] = "\u2659", "\u265F"
	r["b"], r["B"] = "\u2657", "\u265D"
	r["n"], r["N"] = "\u2658", "\u265E"
	r["r"], r["R"] = "\u2656", "\u265C"
	r["q"], r["Q"] = "\u2655", "\u265B"
	r["k"], r["K"] = "\u2654", "\u265A"
	r["."] = "\u00B7"
	// Rows for PGN parsing
	rows := make(map[int][8]int)
	rows[1] = [8]int{18, 17, 16, 15, 14, 13, 12, 11}
	rows[2] = [8]int{28, 27, 26, 25, 24, 23, 22, 21}
	rows[3] = [8]int{38, 37, 36, 35, 34, 33, 32, 31}
	rows[4] = [8]int{48, 47, 46, 45, 44, 43, 42, 41}
	rows[5] = [8]int{58, 57, 56, 55, 54, 53, 52, 51}
	rows[6] = [8]int{68, 67, 66, 65, 64, 63, 62, 61}
	rows[7] = [8]int{78, 77, 76, 75, 74, 73, 72, 71}
	rows[8] = [8]int{88, 87, 86, 85, 84, 83, 82, 81}
	// Regex Patterns
	pgnPattern, _ := regexp.Compile(`([PNBRQK]?[a-h]?[1-8]?)x?([a-h][1-8])([\+\?\!]?)|O(-?O){1,2}`)
	fenPattern, _ := regexp.Compile(`([PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8})\s(w|b)\s([KQkq-]{1,4})\s([a-h][36]|-)\s\d\s([1-9]?[1-9])`)
	return Board{
		board:      b,
		castle:     []byte(`KQkq`),
		pgnMap:     m,
		pieceMap:   p,
		PieceMap:   p,
		pieces:     r,
		rows:       rows,
		toMove:     "w",
		score:      "*",
		moves:      1,
		pgnPattern: pgnPattern,
		fenPattern: fenPattern,
	}
}

// String() returns a string printable board.
// The board will rotate according to whose turn
// it is.
func (b *Board) String() string {
	var printBoard string
	if b.toMove == "w" {
		printBoard = b.StringWhite()
	} else {
		printBoard = b.StringBlack()
	}
	return printBoard
}

// StringWhite returns a string printable board
// from white's perspective.
func (b *Board) StringWhite() string {
	var nums [8]byte // somehow print these?
	nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], nums[6], nums[7] = '1', '2', '3', '4', '5', '6', '7', '8'
	r := make(map[int]bool) // black squares
	r[17], r[15], r[13], r[11], r[28], r[26], r[24], r[22], r[37], r[35], r[33], r[31], r[48], r[46], r[44], r[42], r[57], r[55], r[53], r[51], r[68], r[66], r[64], r[62], r[77], r[75], r[73], r[71], r[88], r[86], r[84], r[82] = false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false
	game := b.board
	p := b.pieces
	var printBoard string
	j := 7
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += string(nums[j]) + ": " + "\n"
			j--
			continue
		} else if (i+1)%10 == 0 {
			continue
		}
		if game[i] == '.' {
			_, ok := r[i]
			if ok { // white square
				printBoard += "|" + " " + "|"
			} else { // black square
				printBoard += "|" + "\u2591" + "|"
			}
		} else {
			printBoard += "|" + p[string(game[i])] + "|"
		}
	}

	printBoard += string(nums[j]) + ": " + "\n"
	printBoard += ":a::b::c::d::e::f::g::h:\n"
	return printBoard
}

// StringBlack rotates for Black perspective.
func (b *Board) StringBlack() string {
	var nums [8]byte // somehow print these?
	nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], nums[6], nums[7] = '1', '2', '3', '4', '5', '6', '7', '8'
	r := make(map[int]bool) // black squares
	r[17], r[15], r[13], r[11], r[28], r[26], r[24], r[22], r[37], r[35], r[33], r[31], r[48], r[46], r[44], r[42], r[57], r[55], r[53], r[51], r[68], r[66], r[64], r[62], r[77], r[75], r[73], r[71], r[88], r[86], r[84], r[82] = false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false
	game := b.board
	p := b.pieces
	var printBoard string
	j := 0
	for i := 11; i < 90; i++ {
		if i%10 == 0 {
			printBoard += string(nums[j]) + ": " + "\n"
			j++
			continue
		} else if (i+1)%10 == 0 {
			continue
		}
		if game[i] == '.' {
			_, ok := r[i]
			if ok { // white square
				printBoard += "|" + " " + "|"
			} else { // black square
				printBoard += "|" + "\u2591" + "|"
			}
		} else {
			printBoard += "|" + p[string(game[i])] + "|"
		}
	}

	printBoard += string(nums[j]) + ": " + "\n"
	printBoard += ":h::g::f::e::d::c::b::a:\n"
	return printBoard
}

// PgnString returns headers and pgn history.
func (b *Board) PgnString() string {
	return b.headers + b.pgn
}

// Position returns string FEN position.
// It also sets the Board.fen attribute
// to the most currect position. (b.fen
// remains empty unil b.Position() is called)
func (b *Board) Position() string {
	pos := ""
	emp := "-"
	zeroTicker := 0
	for i := 88; i > 10; i-- {
		if i%10 == 0 || (i+1)%10 == 0 {
			continue
		}
		// Cycle backwards and tally empty squares
		if b.board[i] == '.' {
			zeroTicker++
		} else if zeroTicker > 0 && b.board[i] != '.' {
			pos += strconv.Itoa(zeroTicker)
			pos += string(b.board[i])
			zeroTicker = 0
		} else {
			pos += string(b.board[i])
		}
		if (i-1)%10 == 0 && i > 10 { // hit edge
			if zeroTicker > 0 {
				pos += strconv.Itoa(zeroTicker)
			}
			zeroTicker = 0
			if i > 11 {
				pos += "/"
			}
		}

	}

	if b.empassant != 0 {
		if b.toMove == "w" {
			emp = b.pieceMap[b.empassant+10]
		} else {
			emp = b.pieceMap[b.empassant-10]
		}
	}
	b.fen = pos + " " + b.toMove + " " + string(b.castle[:4]) + " " + emp + " 0 " + strconv.Itoa(b.moves)
	return b.fen
}

// Stats returns program data of current game
// in map[string]string.
// Todo, replace with exported struct attirbutes.
func (b *Board) Stats() map[string]string {
	_ = b.Position()
	m := make(map[string]string)
	m["turn"] = b.toMove
	m["move"] = strconv.Itoa(b.moves)
	m["castling"] = string(b.castle)
	m["position"] = b.fen
	m["history"] = b.pgn
	m["check"] = strconv.FormatBool(b.Check)
	m["headers"] = b.headers
	m["score"] = b.score
	m["checkmate"] = strconv.FormatBool(b.Checkmate)
	//m["lastthree"] =
	return m
}

// SetHeaders sets pgnHeaders for a pgn export.
func (b *Board) SetHeaders(w, bl string) {
	w = strings.TrimRight(w, "\r\n")
	bl = strings.TrimRight(bl, "\r\n")
	y, m, d := time.Now().Date()
	ye, mo, da := strconv.Itoa(y), strconv.Itoa(int(m)),
		strconv.Itoa(d)
	white := "[White \"" + w + "\"]"
	black := "[Black \"" + bl + "\"]"
	date := "[Date \"" + ye + "." + mo + "." + da + "\"]"
	result := `[Result "*"]`
	b.headers = white + "\n" + black + "\n" + date + "\n" + result + "\n"
}

// Play game in terminal
func main() {
	board := NewBoard()
	PlayGame(board)
}

// PlayGame takes user input and commands.
// See ui/clichess.go for more robust client.
func PlayGame(board Board) {
	var turn string
	welcome := `
********
    go-chess

    /~ |_ _  _ _
    \_|||(/__\_\

`
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(welcome)
	fmt.Print(board.String())

Loop:
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
			case input == "/quit":
				break Loop //os.Exit(1)
			case input == "/new":
				board = NewBoard()
				fmt.Print(board.String())
			case input == "/print":
				fmt.Print(board.String())
			default:
				fmt.Println("Mysterious input")
			}
			continue
		}
		e := board.ParseMove(input)
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Println("\n-------------------")
		// TODO use formats.
		if e != nil {
			fmt.Printf("   [Error: %v]\n", e)
		}
		fmt.Print(board.String())
		if board.Checkmate {
			fmt.Println("****CheckMate!****")
		} else if board.Check {
			fmt.Println("****Check!****")
		}
	}
	fmt.Println("\nGood Game.")
}

// Coordinates prints the int values used
// for Board.Move()
func (b *Board) Coordinates() {
	// TODO Rotate Board
	game := b.board
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i]) + ": "
			continue
		}
		printBoard += "|" + strconv.Itoa(i) + "|"
	}
	printBoard += "\n"
	printBoard += "   :a ::b ::c ::d ::e ::f ::g ::h :\n"
	fmt.Println(printBoard)
}

// isUpper is a wrapper to check if byte in
// Board.board is upper case.
// If Uppercase, it is either white player
// [TODO] or it is empty square.
func (b Board) isUpper(x int) bool {
	return b.board[x] < []byte{0x5a}[0]
}

func (b *Board) cycleHistory(o, d int) {
	b.history[7] = b.history[5]
	b.history[6] = b.history[4]
	b.history[5] = b.history[3]
	b.history[4] = b.history[2]
	b.history[3] = b.history[1]
	b.history[2] = b.history[0]
	//b.history[i-1] = b.history[i-2]
	//fmt.Println(b.history[i-2])
	b.history[1] = d
	b.history[0] = o

}

// CopyBoard takes in a Board pointer and returns
// a copy of it's state, this is for modifying and then
// keeping the originals state intact. One must be careful because
// the values of Board are []byte slices, and these are themselves
// pointers.
func CopyBoard(b *Board) *Board {
	c := *b                        // dereference the pointer
	boardCopy := make([]byte, 120) // []bytes are slices
	castleCopy := make([]byte, 4)
	copy(boardCopy, b.board)
	copy(castleCopy, b.castle)
	c.board = boardCopy
	c.castle = castleCopy
	return &c
}
