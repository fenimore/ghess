package ghess

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ParseStand does pas grande chose mnt.
// TODO: include piece value? Or use different method
// Pgn map takes standar notation and returns coordinate.
// This method is useful for chessboardjs gui
func (b *Board) ParseStand(orig, dest string) error {
	e := b.Move(b.pgnMap[orig], b.pgnMap[dest])
	if e != nil {
		return e
	}
	return nil
}

// ParseMove infers origin and destination
// coordinates from a pgn notation move. Check
// and Check Mate notations will be added automatically.
// TODO: disambiguiation
func (b *Board) ParseMove(move string) error {
	move = strings.TrimRight(move, "\r\n") // prepare for input
	// Variables
	var piece string    // find move piece
	var orig int        // find origin coord of move
	var square string   // find pgnMap key of move
	var attacker string // left of x
	//var precise string // for multiple possibilities
	var target byte // the piece to move, in proper case
	var column int
	var row [8]int

	columns := make(map[string]int)
	columns["a"] = 8
	columns["b"] = 7
	columns["c"] = 6
	columns["d"] = 5
	columns["e"] = 4
	columns["f"] = 3
	columns["g"] = 2
	columns["h"] = 1

	// Status
	isCastle := false
	isWhite := b.toMove == "w"
	isCapture, _ := regexp.MatchString(`x`, move)

	res := b.pgnPattern.FindStringSubmatch(move)
	if res == nil && move != "O-O" && move != "O-O-O" {
		return errors.New("invalid input")
	} else if move == "O-O" || move == "O-O-O" {
		// be nice
		isCastle = true
	}
	// Either is capture or not
	if isCapture {
		attacker = res[1]
		if attacker == strings.ToLower(attacker) {
			piece = "P"
			column = columns[attacker]
		} else { // if  upper case, forcement a piece
			piece = string(attacker[0])
			if len(res[1]) > 1 {
				i, err := strconv.Atoi(string(res[1][1]))
				if err == nil && i != 0 {
					row = b.rows[i]
				} else {
					c := string(res[1][1])
					column = columns[c]
				}
			}
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
			if len(res[1]) > 1 {
				i, err := strconv.Atoi(string(res[1][1]))
				if err == nil && i != 0 {
					row = b.rows[i]
				} else {
					c := string(res[1][1])
					column = columns[c]
				}
			}
			piece = string(res[1][0]) // remove second char
			square = res[2]
		} else {
			return errors.New("Not enough input")
		}
	}
	// the presumed destination
	dest := b.pgnMap[square]
	// The piece will be saved as case sensitive byte
	if b.toMove == "b" {
		if piece != "" {
			target = []byte(strings.ToLower(piece))[0]
		} else {
			return errors.New("Invalid Input")
		}
	} else {
		if piece != "" {
			target = []byte(piece)[0]
		} else {
			return errors.New("Invalid Input")
		}
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
		// Disambiguate
		firstDisambig := column != 0 && possibilities[0]%10 == column
		secondDisambig := column != 0 && possibilities[1]%10 == column
		firstPoss := b.board[possibilities[0]] == target
		secondPoss := b.board[possibilities[1]] == target
		if firstPoss && firstDisambig && isCapture {
			orig = possibilities[0]
		} else if secondPoss && secondDisambig && isCapture {
			orig = possibilities[1]
		} else if firstPoss {
			orig = possibilities[0]
		} else if secondPoss {
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
			if column != 0 { // Disambiguate
				disambig := possibility%10 == column
				if b.board[possibility] == target && disambig {
					orig = possibility
					break
				}
			} else if row[0] != 0 { // Disambiguate
				for _, r := range row {
					disambig := b.board[possibility] == target &&
						b.board[possibility] == byte(r)
					if disambig {
						orig = possibility
						break
					}
				}
			} else {
				if b.board[possibility] == target {
					orig = possibility
					break
				}
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
	LoopRook:
		for _, possibility := range possibilities {
			if column != 0 { // Disambiguate
				disambig := possibility%10 == column
				if b.board[possibility] == target && disambig {
					orig = possibility
					err := b.validRook(orig, dest)
					if err != nil {
						continue
					}
					break LoopRook
				}
			} else if row[0] != 0 { // Disambiguate
				for _, r := range row {
					disambig := b.board[possibility] == target &&
						possibility == r
					if disambig {
						orig = possibility
						err := b.validRook(orig, dest)
						if err != nil {
							continue
						}
						break LoopRook
					}
				}
			} else {
				if b.board[possibility] == target {
					orig = possibility
					err := b.validRook(orig, dest)
					if err != nil {
						continue
					}
					break LoopRook
				}
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
		return errors.New("Not the proper capture syntax")

	}
	if orig != 0 && dest != 0 {
		err := b.Move(orig, dest)
		if err == nil {
			matchCheck, _ := regexp.MatchString(`+`, move)
			matchMate, _ := regexp.MatchString(`#`, move)
			// Update pgn History
			if b.toMove == "b" {
				b.pgn += strconv.Itoa(b.moves) + ". "
			}
			b.pgn += (move)
			if b.Checkmate && matchMate {
				b.pgn += "# "
			} else if b.Check && matchCheck {

				b.pgn += "+ "
			} else {
				b.pgn += " " // add space
			}
		}
		return err
	}
	return errors.New("No such move")
}

// LoadPgn reads a pgn match.
// TODO: ignore header strings eg [White].
func (b *Board) LoadPgn(match string) error {
	// Does this already ignore?
	result := b.pgnPattern.FindAllString(match, -1)
	for _, val := range result {
		err := b.ParseMove(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadFen parses FEN string and update Board.board.
func (b *Board) LoadFen(fen string) error {
	// Treat fen input
	fen = strings.TrimRight(fen, "\r\n")
	matches := b.fenPattern.MatchString(fen)
	if !matches {
		return errors.New("Invalid FEN")
	}
	res := b.fenPattern.FindStringSubmatch(fen)
	posCount := 88 // First position to fill
	//
	for _, val := range res[1] { // res[1] {
		// First, make sure it's a relevant position
		if (posCount % 10) == 0 { // || (posCount+1)%10 == 0 {
			posCount -= 2
		}
		// Check if there are Empty Squares
		num, e := strconv.Atoi(string(val))
		if e == nil {
			for j := 0; j < num; j++ {
				b.board[posCount] = '.'
				posCount--
			}
			continue
		}
		if val == '/' {
			continue
		}
		// Is regular Piece
		b.board[posCount] = byte(val)
		posCount--
	}
	// Update Board value
	// Castle
	b.castle = []byte("----")
	for _, v := range res[3] {
		switch {
		case v == 'K':
			b.castle[0] = 'K'
			break
		case v == 'Q':
			b.castle[1] = 'Q'
			break
		case v == 'k':
			b.castle[2] = 'k'
			break
		case v == 'q':
			b.castle[3] = 'q'
			break
		}
	}
	// Empassant
	if res[4] != "-" {
		emp := b.pgnMap[res[4]]
		if emp < 40 {
			// white pawn
		} else if emp > 60 {
		} else {
			return errors.New("Invalid FEN empassant")
		}
	} else {
		b.empassant = 0
	}
	// Turn
	turns, _ := strconv.Atoi(res[5])
	b.moves = turns
	b.fen = fen
	b.toMove = res[2]
	b.Check = b.isPlayerInCheck()
	return nil
}
