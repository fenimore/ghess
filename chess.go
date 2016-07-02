/*
Search and Evaluation


*/
package main

import (
	"errors"
	"fmt"
	//"time"
)

//var Board [128]byte

/*
    ' rnbqkbnr\n'  #  20 - 29
    ' pppppppp\n'  #  30 - 39
    ' ........\n'  #  40 - 49
    ' ........\n'  #  50 - 59
    ' ........\n'  #  60 - 69
    ' ........\n'  #  70 - 79
    ' PPPPPPPP\n'  #  80 - 89
    ' RNBQKBNR\n'  #  90 - 99


*/
type Board struct {
	board []byte
	A1, H1, A8, H8 int
}

func (b Board) String() string {
	var printBoard string
	
	fmt.Print(string(b.board))
	//* Good for 64 byte board
	for idx, _ := range b.board {
		if idx < 10 {
			fmt.Print(":00", idx, "  ")
		}
		if idx < 100 && idx > 10{
			fmt.Print( ":0", idx, "  ")
		}
		if idx > 99 {
			fmt.Print(":", idx, "  ")
		}
		sq := idx 
		if sq % 10 == 0 && sq != 0{
			fmt.Print("\n")
		}
	}
			
	return printBoard
}

func (b Board) Set(token byte, x, y int) error {
	idx, err := b.squareAt(x, y)
	if err != nil {
		return err
	}
	b.board[idx] = token
	return nil
}

func (b Board) squareAt(x, y int) (int, error) {
	woff := 16
	foff := (y*2+1)*woff + x*2 + 1

	if foff > len(b.board){
		return 0, errors.New("out of range")
	}
	return  (y*2+1)*woff + x*2 + 1, nil
}

func NewBoard() Board {
	b := make([]byte, 120)

	/*for i := 0; i < 8; i +=2 {
		row0 := i * 8
		row1 := (i + 1) * 8
		for j := 0; j < 8; j += 2 {
			b[row0+j], b[row0+j+1] = 'b', 'w'
			if row1+j+1 <= 64 {
				b[row1+j], b[row1+j+1] = 'w', 'b'
			}
		}
	}*/
	fmt.Println("initializing board")
	//fmt.Println("At 20-29", string(b[28]))
	return Board{
		board: b,
		A1: 91,
		A8: 21,
		H1: 98,
		H8: 28,
	}
}


func main() {
//	Board[0] = 70
//	Board[1] = 
	///	fmt.Print(Board)
	board := NewBoard()
	//fmt.Print(board.String())
	//fmt.Println(string(board.board[1:8]))
	board.String()
	fmt.Println("At 20-29", string(board.board[31]))
}
