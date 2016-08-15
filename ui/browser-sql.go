/*
TODO:
Create templates:
Index with ids,
get set player data
templates with forms
add javascript chessboard
*/
package main

import (
	"fmt"
	"github.com/polypmer/ghess"
	"net/http"
)

type ChessHandler struct {
	g *Board
}

// boardHandler for playing game
// Takes url param pgn move
func (h *ChessHandler) playGameHandler(w http.ResponseWriter,
	r *http.Request) {
	game := ghess.NewBoard()
	fmt.Fprintln(w, game.String())
	move := r.URL.Path[1:]
	e := game.ParseMove(move)
	if e != nil {
		fmt.Fprintln(w, e.Error())
	}
	fmt.Fprintln(w, game.String())
}

func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	// Create New Game
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	h := new(ChessHandler)

	// Server Part
	http.HandleFunc("/", h.playGameHandler)
	http.HandleFunc("/new", h.newGameHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
