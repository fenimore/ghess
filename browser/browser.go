package main

import (
	"fmt"
	"github.com/polypmer/ghess"
	"net/http"
)

func boardHandler(w http.ResponseWriter, r *http.Request) {
	game := ghess.NewBoard()
	fmt.Fprintf(w, game.String())
	move := r.URL.Path[1:]
	game.ParseMove(move)
	fmt.Fprintf(w, game.String())
}

func main() {
	// Listen
	http.HandleFunc("/", boardHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
