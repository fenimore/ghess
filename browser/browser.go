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
	e := game.ParseMove(move)
	if e != nil {
		fmt.Fprintf(w, e.Error()+"\n")
	}
	fmt.Fprintf(w, game.String())
}

func main() {
	// Listen
	// So ahndFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	http.HandleFunc("/", boardHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
